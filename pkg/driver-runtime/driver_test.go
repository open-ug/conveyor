package driverruntime_test

import (
	"testing"

	driverruntime "github.com/open-ug/conveyor/pkg/driver-runtime"
	"github.com/open-ug/conveyor/pkg/driver-runtime/log"
	"github.com/open-ug/conveyor/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDriver_Validate(t *testing.T) {
	tests := []struct {
		name    string
		driver  driverruntime.Driver
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid driver",
			driver: driverruntime.Driver{
				Reconcile: func(message, event, runID string, logger *log.DriverLogger) types.DriverResult {
					return types.DriverResult{}
				},
				Name:      "test-driver",
				Resources: []string{"pods", "services"},
			},
			wantErr: false,
		},
		{
			name: "missing reconcile",
			driver: driverruntime.Driver{
				Reconcile: nil,
				Name:      "test-driver",
				Resources: []string{"pods"},
			},
			wantErr: true,
			errMsg:  "driver reconcile function is not set",
		},
		{
			name: "missing name",
			driver: driverruntime.Driver{
				Reconcile: func(message, event, runID string, logger *log.DriverLogger) types.DriverResult {
					return types.DriverResult{}
				},
				Name:      "",
				Resources: []string{"pods"},
			},
			wantErr: true,
			errMsg:  "driver name is not set",
		},
		{
			name: "missing resources",
			driver: driverruntime.Driver{
				Reconcile: func(message, event, runID string, logger *log.DriverLogger) types.DriverResult {
					return types.DriverResult{}
				},
				Name:      "test-driver",
				Resources: []string{},
			},
			wantErr: true,
			errMsg:  "driver resources are not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.driver.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
