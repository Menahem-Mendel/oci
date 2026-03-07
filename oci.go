package oci

import (
	"fmt"
	"oci/driver"
	"sort"
	"strings"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

// Register registers a runtime driver by name.
// It panics if the driver is nil or if the name is already registered.
func Register(name string, drv driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if drv == nil {
		panic("oci: Register nil driver")
	}
	if name == "" {
		panic("oci: Register empty driver name")
	}
	if _, dup := drivers[name]; dup {
		panic("oci: Register called twice for driver " + name)
	}
	drivers[name] = drv
}

// Drivers returns registered driver names in sorted order.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()

	names := make([]string, 0, len(drivers))
	for name := range drivers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// NewRuntime creates a runtime handle from a registered driver name.
// It does not open the driver connection; call Connect or Open.
func NewRuntime(driverName string) (*Runtime, error) {
	driverName = strings.TrimSpace(driverName)
	if driverName == "" {
		return nil, fmt.Errorf("%w: driver name is required", ErrInvalidArgument)
	}

	driversMu.RLock()
	drv, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnregisteredDriver, driverName)
	}

	return &Runtime{driverName: driverName, drv: drv}, nil
}

// Open creates and connects a runtime in one call.
func Open(driverName, dsn string) (*Runtime, error) {
	rt, err := NewRuntime(driverName)
	if err != nil {
		return nil, err
	}
	if err := rt.Connect(dsn); err != nil {
		return nil, err
	}
	return rt, nil
}
