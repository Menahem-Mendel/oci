package podman

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"oci"
	"oci/driver"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultHTTPTimeout = 60 * time.Second

func init() {
	oci.Register("podman", &Driver{})
}

var (
	_ driver.Driver    = (*Driver)(nil)
	_ driver.Conn      = (*Conn)(nil)
	_ driver.Pinger    = (*Conn)(nil)
	_ driver.Puller    = (*Conn)(nil)
	_ driver.Inspector = (*Conn)(nil)
	_ driver.Lister    = (*Conn)(nil)
	_ driver.Remover   = (*Conn)(nil)
	_ driver.Creator   = (*Conn)(nil)
	_ driver.Starter   = (*Conn)(nil)
	_ driver.Stopper   = (*Conn)(nil)
	_ driver.Execer    = (*Conn)(nil)
)

// Driver implements the oci/driver.Driver interface for the Podman API.
// It uses Podman's Docker-compatible HTTP API.
type Driver struct {
	HTTPClient *http.Client
}

func (d *Driver) Open(dsn string) (driver.Conn, error) {
	u, err := url.Parse(strings.TrimSpace(dsn))
	if err != nil {
		return nil, fmt.Errorf("podman: parse dsn: %w", err)
	}

	baseURL := ""
	client := d.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: defaultHTTPTimeout}
	}

	switch u.Scheme {
	case "unix":
		socket := u.Path
		if socket == "" {
			socket = u.Host
		}
		if socket == "" {
			return nil, fmt.Errorf("podman: unix dsn requires socket path")
		}

		transport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, "unix", socket)
			},
		}
		client = &http.Client{Transport: transport, Timeout: client.Timeout}
		if client.Timeout == 0 {
			client.Timeout = defaultHTTPTimeout
		}
		baseURL = "http://d"

	case "tcp":
		if u.Host == "" {
			return nil, fmt.Errorf("podman: tcp dsn requires host")
		}
		baseURL = "http://" + u.Host

	case "http", "https":
		if u.Host == "" {
			return nil, fmt.Errorf("podman: %s dsn requires host", u.Scheme)
		}
		baseURL = u.Scheme + "://" + u.Host

	default:
		return nil, fmt.Errorf("podman: unsupported dsn scheme %q", u.Scheme)
	}

	return &Conn{client: client, baseURL: baseURL}, nil
}

type Conn struct {
	client  *http.Client
	baseURL string

	mu     sync.RWMutex
	closed bool
}

func (c *Conn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	if c.client != nil {
		c.client.CloseIdleConnections()
	}
	return nil
}

func (c *Conn) Ping(ctx context.Context) error {
	resp, err := c.do(ctx, http.MethodGet, "/_ping", nil, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return apiError(resp)
}

func (c *Conn) Pull(ctx context.Context, reference string, options map[string]string) (string, error) {
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return "", fmt.Errorf("podman: image reference is required")
	}

	query := url.Values{}
	query.Set("fromImage", reference)

	if platform := strings.TrimSpace(options[oci.OptionPlatform]); platform != "" {
		query.Set("platform", platform)
	}
	if parseBool(options[oci.OptionAllTags]) {
		query.Set("allTags", "true")
	}
	if parseBool(options[oci.OptionInsecureSkipTLS]) {
		query.Set("tlsVerify", "false")
	}

	headers := map[string]string{}
	username := options[oci.OptionUsername]
	password := options[oci.OptionPassword]
	if username != "" || password != "" {
		authRaw, err := json.Marshal(map[string]string{
			"username": username,
			"password": password,
		})
		if err != nil {
			return "", err
		}
		headers["X-Registry-Auth"] = base64.StdEncoding.EncodeToString(authRaw)
	}

	resp, err := c.do(ctx, http.MethodPost, "/images/create", query, nil, headers)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", apiError(resp)
	}
	_, _ = io.Copy(io.Discard, resp.Body)

	doc, err := c.Inspect(ctx, driver.KindImage, reference)
	if err != nil {
		return "", err
	}
	if id, _ := doc["id"].(string); strings.TrimSpace(id) != "" {
		return id, nil
	}
	return reference, nil
}

func (c *Conn) Inspect(ctx context.Context, kind string, idOrRef string) (map[string]any, error) {
	switch kind {
	case driver.KindImage:
		return c.inspectImage(ctx, idOrRef)
	case driver.KindContainer:
		return c.inspectContainer(ctx, idOrRef)
	default:
		return nil, fmt.Errorf("podman: inspect kind %q: %w", kind, driver.ErrNotSupported)
	}
}

func (c *Conn) List(ctx context.Context, kind string, filters map[string]string) ([]map[string]any, error) {
	switch kind {
	case driver.KindImage:
		return c.listImages(ctx, filters)
	case driver.KindContainer:
		return c.listContainers(ctx, filters)
	default:
		return nil, fmt.Errorf("podman: list kind %q: %w", kind, driver.ErrNotSupported)
	}
}

func (c *Conn) Remove(ctx context.Context, kind string, idOrRef string, force bool) error {
	switch kind {
	case driver.KindImage:
		query := url.Values{}
		if force {
			query.Set("force", "1")
		}
		resp, err := c.do(ctx, http.MethodDelete, "/images/"+url.PathEscape(idOrRef), query, nil, nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		return apiError(resp)

	case driver.KindContainer:
		query := url.Values{}
		if force {
			query.Set("force", "1")
		}
		resp, err := c.do(ctx, http.MethodDelete, "/containers/"+url.PathEscape(idOrRef), query, nil, nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		return apiError(resp)

	default:
		return fmt.Errorf("podman: remove kind %q: %w", kind, driver.ErrNotSupported)
	}
}

func (c *Conn) Create(ctx context.Context, kind string, spec map[string]any) (string, error) {
	switch kind {
	case driver.KindContainer:
		return c.createContainer(ctx, spec)
	default:
		return "", fmt.Errorf("podman: create kind %q: %w", kind, driver.ErrNotSupported)
	}
}

func (c *Conn) Start(ctx context.Context, kind string, id string) error {
	switch kind {
	case driver.KindContainer:
		resp, err := c.do(ctx, http.MethodPost, "/containers/"+url.PathEscape(id)+"/start", nil, nil, nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		return apiError(resp)
	default:
		return fmt.Errorf("podman: start kind %q: %w", kind, driver.ErrNotSupported)
	}
}

func (c *Conn) Stop(ctx context.Context, kind string, id string, timeoutSeconds int) error {
	switch kind {
	case driver.KindContainer:
		query := url.Values{}
		query.Set("t", strconv.Itoa(timeoutSeconds))
		resp, err := c.do(ctx, http.MethodPost, "/containers/"+url.PathEscape(id)+"/stop", query, nil, nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		return apiError(resp)
	default:
		return fmt.Errorf("podman: stop kind %q: %w", kind, driver.ErrNotSupported)
	}
}

func (c *Conn) Exec(ctx context.Context, id string, command []string, options map[string]any) (map[string]any, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("podman: container id is required")
	}
	if len(command) == 0 {
		return nil, fmt.Errorf("podman: command is required")
	}

	env := toEnvMap(options[oci.ExecOptionEnv])
	tty := toBool(options[oci.ExecOptionTTY])
	stdin := toBytes(options[oci.ExecOptionStdin])

	createBody, err := json.Marshal(map[string]any{
		"AttachStdin":  len(stdin) > 0,
		"AttachStdout": true,
		"AttachStderr": true,
		"Tty":          tty,
		"Cmd":          command,
		"Env":          envSlice(env),
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, http.MethodPost, "/containers/"+url.PathEscape(id)+"/exec", nil, bytes.NewReader(createBody), map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return nil, apiError(resp)
	}
	var created struct {
		ID string `json:"Id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("podman: decode exec create: %w", err)
	}
	resp.Body.Close()

	startBody, err := json.Marshal(map[string]any{
		"Detach": false,
		"Tty":    tty,
	})
	if err != nil {
		return nil, err
	}

	resp, err = c.do(ctx, http.MethodPost, "/exec/"+url.PathEscape(created.ID)+"/start", nil, bytes.NewReader(startBody), map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return nil, apiError(resp)
	}
	rawOutput, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	stdout := rawOutput
	stderr := []byte(nil)
	if !tty {
		stdout, stderr = demuxDockerStream(rawOutput)
	}

	resp, err = c.do(ctx, http.MethodGet, "/exec/"+url.PathEscape(created.ID)+"/json", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apiError(resp)
	}

	var inspect struct {
		ExitCode int `json:"ExitCode"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&inspect); err != nil {
		return nil, fmt.Errorf("podman: decode exec inspect: %w", err)
	}

	return map[string]any{
		oci.ExecResultExitCode: inspect.ExitCode,
		oci.ExecResultStdout:   stdout,
		oci.ExecResultStderr:   stderr,
	}, nil
}

func (c *Conn) inspectImage(ctx context.Context, idOrRef string) (map[string]any, error) {
	resp, err := c.do(ctx, http.MethodGet, "/images/"+url.PathEscape(idOrRef)+"/json", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apiError(resp)
	}

	var payload struct {
		ID          string   `json:"Id"`
		RepoTags    []string `json:"RepoTags"`
		RepoDigests []string `json:"RepoDigests"`
		Size        int64    `json:"Size"`
		Created     string   `json:"Created"`
		Config      struct {
			Labels map[string]string `json:"Labels"`
		} `json:"Config"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("podman: decode image inspect: %w", err)
	}

	doc := map[string]any{
		"kind":       driver.KindImage,
		"id":         payload.ID,
		"reference":  firstNonEmpty(payload.RepoTags, idOrRef),
		"digest":     firstDigest(payload.RepoDigests),
		"size_bytes": payload.Size,
		"labels":     payload.Config.Labels,
	}
	if created, err := time.Parse(time.RFC3339Nano, payload.Created); err == nil {
		doc["created_at"] = created.UTC().Format(time.RFC3339Nano)
	}
	return doc, nil
}

func (c *Conn) listImages(ctx context.Context, _ map[string]string) ([]map[string]any, error) {
	resp, err := c.do(ctx, http.MethodGet, "/images/json", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apiError(resp)
	}

	var payload []struct {
		ID          string            `json:"Id"`
		RepoTags    []string          `json:"RepoTags"`
		RepoDigests []string          `json:"RepoDigests"`
		Size        int64             `json:"Size"`
		Created     int64             `json:"Created"`
		Labels      map[string]string `json:"Labels"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("podman: decode images list: %w", err)
	}

	out := make([]map[string]any, 0, len(payload))
	for _, item := range payload {
		out = append(out, map[string]any{
			"kind":       driver.KindImage,
			"id":         item.ID,
			"reference":  firstNonEmpty(item.RepoTags, item.ID),
			"digest":     firstDigest(item.RepoDigests),
			"size_bytes": item.Size,
			"created_at": time.Unix(item.Created, 0).UTC().Format(time.RFC3339Nano),
			"labels":     item.Labels,
		})
	}
	return out, nil
}

func (c *Conn) createContainer(ctx context.Context, spec map[string]any) (string, error) {
	image := strings.TrimSpace(toString(spec[oci.SpecImage]))
	if image == "" {
		return "", fmt.Errorf("podman: container image is required in spec[%q]", oci.SpecImage)
	}

	query := url.Values{}
	if name := strings.TrimSpace(toString(spec[oci.SpecName])); name != "" {
		query.Set("name", name)
	}

	body := map[string]any{
		"Image": image,
	}

	if cmd := toStringSlice(spec[oci.SpecCommand]); len(cmd) > 0 {
		body["Cmd"] = cmd
	}
	if env := toEnvMap(spec[oci.SpecEnv]); len(env) > 0 {
		body["Env"] = envSlice(env)
	}
	if labels := toStringMap(spec[oci.SpecLabels]); len(labels) > 0 {
		body["Labels"] = labels
	}
	if wd := strings.TrimSpace(toString(spec[oci.SpecWorkingDir])); wd != "" {
		body["WorkingDir"] = wd
	}
	if networkMode := strings.TrimSpace(toString(spec[oci.SpecNetworkMode])); networkMode != "" {
		body["HostConfig"] = map[string]any{"NetworkMode": networkMode}
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	resp, err := c.do(ctx, http.MethodPost, "/containers/create", query, bytes.NewReader(raw), map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", apiError(resp)
	}

	var created struct {
		ID string `json:"Id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return "", fmt.Errorf("podman: decode container create: %w", err)
	}
	return created.ID, nil
}

func (c *Conn) inspectContainer(ctx context.Context, id string) (map[string]any, error) {
	resp, err := c.do(ctx, http.MethodGet, "/containers/"+url.PathEscape(id)+"/json", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apiError(resp)
	}

	var payload struct {
		ID      string `json:"Id"`
		Name    string `json:"Name"`
		Created string `json:"Created"`
		Config  struct {
			Image  string            `json:"Image"`
			Cmd    []string          `json:"Cmd"`
			Labels map[string]string `json:"Labels"`
		} `json:"Config"`
		State struct {
			Status string `json:"Status"`
		} `json:"State"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("podman: decode container inspect: %w", err)
	}

	doc := map[string]any{
		"kind":       driver.KindContainer,
		"id":         payload.ID,
		"name":       strings.TrimPrefix(payload.Name, "/"),
		"image":      payload.Config.Image,
		"command":    payload.Config.Cmd,
		"state":      payload.State.Status,
		"status":     payload.State.Status,
		"labels":     payload.Config.Labels,
		"created_at": payload.Created,
	}
	if created, err := time.Parse(time.RFC3339Nano, payload.Created); err == nil {
		doc["created_at"] = created.UTC().Format(time.RFC3339Nano)
	}
	return doc, nil
}

func (c *Conn) listContainers(ctx context.Context, filters map[string]string) ([]map[string]any, error) {
	query := url.Values{}
	query.Set("all", "1")
	for k, v := range filters {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			continue
		}
		query.Set(k, v)
	}

	resp, err := c.do(ctx, http.MethodGet, "/containers/json", query, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apiError(resp)
	}

	var payload []struct {
		ID      string            `json:"Id"`
		Names   []string          `json:"Names"`
		Image   string            `json:"Image"`
		Command string            `json:"Command"`
		State   string            `json:"State"`
		Status  string            `json:"Status"`
		Created int64             `json:"Created"`
		Labels  map[string]string `json:"Labels"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("podman: decode containers list: %w", err)
	}

	out := make([]map[string]any, 0, len(payload))
	for _, item := range payload {
		name := ""
		if len(item.Names) > 0 {
			name = strings.TrimPrefix(item.Names[0], "/")
		}
		cmd := []string{}
		if strings.TrimSpace(item.Command) != "" {
			cmd = []string{item.Command}
		}
		out = append(out, map[string]any{
			"kind":       driver.KindContainer,
			"id":         item.ID,
			"name":       name,
			"image":      item.Image,
			"command":    cmd,
			"state":      item.State,
			"status":     item.Status,
			"created_at": time.Unix(item.Created, 0).UTC().Format(time.RFC3339Nano),
			"labels":     item.Labels,
		})
	}
	return out, nil
}

func (c *Conn) do(ctx context.Context, method, endpoint string, query url.Values, body io.Reader, headers map[string]string) (*http.Response, error) {
	c.mu.RLock()
	closed := c.closed
	c.mu.RUnlock()
	if closed {
		return nil, fmt.Errorf("podman: connection closed")
	}

	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("podman: parse base url: %w", err)
	}
	u.Path = endpoint
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("podman: %s %s: %w", method, endpoint, err)
	}
	return resp, nil
}

func apiError(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("podman: nil response")
	}

	payload, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	_ = resp.Body.Close()

	var problem struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	_ = json.Unmarshal(payload, &problem)

	msg := strings.TrimSpace(problem.Message)
	if msg == "" {
		msg = strings.TrimSpace(problem.Error)
	}
	if msg == "" {
		msg = strings.TrimSpace(string(payload))
	}
	if msg == "" {
		msg = resp.Status
	}
	return fmt.Errorf("podman: api %s %s failed: %s", resp.Request.Method, resp.Request.URL.Path, msg)
}

func parseBool(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	return s == "1" || s == "true" || s == "yes"
}

func firstNonEmpty(values []string, fallback string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return fallback
}

func firstDigest(values []string) string {
	for _, v := range values {
		if at := strings.LastIndex(v, "@"); at >= 0 && at < len(v)-1 {
			return v[at+1:]
		}
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func envSlice(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, k+"="+env[k])
	}
	return out
}

func demuxDockerStream(raw []byte) ([]byte, []byte) {
	if len(raw) < 8 {
		return raw, nil
	}

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	i := 0
	for i+8 <= len(raw) {
		stream := raw[i]
		size := int(binary.BigEndian.Uint32(raw[i+4 : i+8]))
		i += 8
		if i+size > len(raw) {
			return raw, nil
		}
		chunk := raw[i : i+size]
		i += size

		switch stream {
		case 1:
			_, _ = stdout.Write(chunk)
		case 2:
			_, _ = stderr.Write(chunk)
		default:
			_, _ = stdout.Write(chunk)
		}
	}
	if i != len(raw) {
		return raw, nil
	}
	return stdout.Bytes(), stderr.Bytes()
}

func toString(v any) string {
	s, _ := v.(string)
	return s
}

func toBool(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return parseBool(t)
	default:
		return false
	}
}

func toBytes(v any) []byte {
	switch t := v.(type) {
	case []byte:
		return append([]byte(nil), t...)
	case string:
		return []byte(t)
	default:
		return nil
	}
}

func toStringSlice(v any) []string {
	switch t := v.(type) {
	case []string:
		return append([]string(nil), t...)
	case []any:
		out := make([]string, 0, len(t))
		for _, item := range t {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

func toStringMap(v any) map[string]string {
	switch t := v.(type) {
	case map[string]string:
		out := make(map[string]string, len(t))
		for k, v := range t {
			out[k] = v
		}
		return out
	case map[string]any:
		out := make(map[string]string, len(t))
		for k, item := range t {
			if s, ok := item.(string); ok {
				out[k] = s
			}
		}
		return out
	default:
		return nil
	}
}

func toEnvMap(v any) map[string]string {
	return toStringMap(v)
}
