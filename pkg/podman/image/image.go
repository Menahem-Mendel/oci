package image

// type Service struct {
// 	conn *bindings.Connection
// }

// func (s *Service) ServeOCI(ctx context.Context, f func(ctx context.Context)) error {
// 	f(ctx)
// 	return nil
// }

// // Service -> HandleImage ->

// // image.Pull -> Service.ServeOCI(puller) -> puller() {puller.Pull}

// // image.Pull -> podman.Pull -> oci.Pull

// type Puller struct {
// 	// conn *bindings.Connection

// 	// AllTags can be specified to pull all tags of an image. Note
// 	// that this only works if the image does not include a tag.
// 	AllTags *bool
// 	// Arch will overwrite the local architecture for image pulls.
// 	Arch *string
// 	// Authfile is the path to the authentication file. Ignored for remote
// 	// calls.
// 	Authfile *string
// 	// OS will overwrite the local operating system (OS) for image
// 	// pulls.
// 	OS *string
// 	// Policy is the pull policy. Supported values are "missing", "never",
// 	// "newer", "always". An empty string defaults to "always".
// 	Policy *string
// 	// Password for authenticating against the registry.
// 	Password *string `schema:"-"`
// 	// ProgressWriter is a writer where pull progress are sent.
// 	ProgressWriter *io.Writer `schema:"-"`
// 	// Quiet can be specified to suppress pull progress when pulling.  Ignored
// 	// for remote calls.
// 	Quiet *bool
// 	// SkipTLSVerify to skip HTTPS and certificate verification.
// 	SkipTLSVerify *bool `schema:"-"`
// 	// Username for authenticating against the registry.
// 	Username *string `schema:"-"`
// 	// Variant will overwrite the local variant for image pulls.
// 	Variant *string
// }

// func Apply[T PullOption](drv driver.Driver)

// func (p *Puller) WithAllTags(opt driver.Option) error {
// 	opt(p.AllTags)
// }

// func (p *Puller) Options(opts ...driver.Option) error {

// }

// func (p *Puller) ServeOCI(ctx context.Context, args ...any) error {
// 	oci.Pull(ctx, p, args...)
// }

// func (s *Puller) Pull(ctx context.Context, ref string) (string, error) {
// 	opts := images.PullOptions{}

// 	if s.conn == nil {
// 		return "", fmt.Errorf("podman: %w", "ErrNoConnection")
// 	}

// 	params, err := opts.ToParams()
// 	if err != nil {
// 		return "", err
// 	}
// 	params.Set("reference", ref)

// 	// SkipTLSVerify is special.  It's not being serialized by ToParams()
// 	// because we need to flip the boolean.
// 	if opts.SkipTLSVerify != nil {
// 		params.Set("tlsVerify", strconv.FormatBool(!opts.GetSkipTLSVerify()))
// 	}

// 	header, err := auth.MakeXRegistryAuthHeader(&types.SystemContext{AuthFilePath: opts.GetAuthfile()}, opts.GetUsername(), opts.GetPassword())
// 	if err != nil {
// 		return "", err
// 	}

// 	response, err := s.conn.DoRequest(ctx, nil, http.MethodPost, "/images/pull", params, header)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer response.Body.Close()

// 	if !response.IsSuccess() {
// 		return "", response.Process(err)
// 	}

// 	dec := json.NewDecoder(response.Body)
// 	var pullErrors []error

// LOOP:
// 	for {
// 		var report entities.ImagePullReport
// 		if err := dec.Decode(&report); err != nil {
// 			if errors.Is(err, io.EOF) {
// 				break
// 			}
// 			report.Error = err.Error() + "\n"
// 		}

// 		select {
// 		case <-response.Request.Context().Done():
// 			break LOOP
// 		default:
// 			// non-blocking select
// 		}

// 		switch {
// 		case report.Stream != "":
// 		case report.Error != "":
// 			pullErrors = append(pullErrors, errors.New(report.Error))
// 		case report.ID != "":
// 		default:
// 			return "", fmt.Errorf("failed to parse pull results stream, unexpected input: %v", report)
// 		}
// 	}

// 	return "", errorhandling.JoinErrors(pullErrors)
// }
