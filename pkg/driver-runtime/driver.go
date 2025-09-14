package driverruntime

import (
	"fmt"

	"github.com/open-ug/conveyor/pkg/driver-runtime/log"
	"github.com/open-ug/conveyor/pkg/types"
)

type Driver struct {
	// The driver is responsible for managing the driver
	Reconcile func(message string, event string, runID string, logger *log.DriverLogger) types.DriverResult

	Name string

	Resources []string
}

// validate the driver
func (d *Driver) Validate() error {
	if d.Reconcile == nil {
		return fmt.Errorf("driver reconcile function is not set")
	}
	if d.Name == "" {
		return fmt.Errorf("driver name is not set")
	}

	if len(d.Resources) == 0 {
		return fmt.Errorf("driver resources are not set")
	}

	return nil
}
