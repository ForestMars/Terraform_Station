package TerraformStation

import (
	"context"
)

// TerraformStationService defines the interface for Terraform operations
type TerraformStationService interface {
	// Basic command execution
	TFCommand(ctx context.Context, input *TFCommandInput) (*TFCommandResult, error)
	
	// Specific Terraform operations
	TFPlan(ctx context.Context, input *TFCommandInput) (*TFPlanResult, error)
	TFApply(ctx context.Context, input *TFCommandInput) (*TFApplyResult, error)
	TFInit(ctx context.Context, input *TFCommandInput) (*TFCommandResult, error)
	TFValidate(ctx context.Context, input *TFCommandInput) (*TFCommandResult, error)
	TFState(ctx context.Context, input *TFCommandInput) (*TFStateInfo, error)
	
	// Utility methods
	GetConfig() *Config
	SetWorkingDirectory(dir string) error
	ValidateWorkingDirectory(dir string) error
}
