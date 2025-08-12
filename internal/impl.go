package internal

import (
	"context"
	"strings"
	"time"

	"github.com/go-gorm/gorm"
	"github.com/ForestMars/TerraformStation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TerraformStationImpl struct {
	db             *gorm.DB
	cfg            *TerraformStation.Config
	executor       *TerraformStation.OpenTofuExecutor
	workingDir     string
}

func New(db *gorm.DB, cfg *TerraformStation.Config) (*TerraformStationImpl, error) {
	if cfg == nil {
		return nil, TerraformStation.NewInvalidInputError("configuration cannot be nil")
	}

	if db == nil {
		return nil, TerraformStation.NewInvalidInputError("database connection cannot be nil")
	}

	// Create opentofu executor
	executor := TerraformStation.NewOpenTofuExecutor(cfg.OpenTofuPath, cfg.Timeout)

	impl := &TerraformStationImpl{
		db:         db,
		cfg:        cfg,
		executor:   executor,
		workingDir: cfg.WorkingDirectory,
	}

	// Validate working directory
	if err := impl.ValidateWorkingDirectory(cfg.WorkingDirectory); err != nil {
		return nil, err
	}

	return impl, nil
}

// TFCommand executes a generic OpenTofu command
func (impl *TerraformStationImpl) TFCommand(ctx context.Context, input *TerraformStation.TFCommandInput) (*TerraformStation.TFCommandResult, error) {
	// Validate input
	if err := TerraformStation.ValidateTFCommandInput(input); err != nil {
		return nil, err
	}

	// Use working directory from input or fallback to configured one
	workingDir := impl.workingDir
	if input.WorkingDirectory != "" {
		workingDir = input.WorkingDirectory
	}

	// Build command arguments
	args := TerraformStation.BuildOpenTofuArgs(input.Command, input)

	// Execute command
	output, err := impl.executor.Execute(ctx, workingDir, args...)
	
	// Create result
	result := &TerraformStation.TFCommandResult{
		CommandId:   TerraformStation.GenerateCommandID(),
		ExecutedAt:  timestamppb.Now(),
		Result:      output,
	}

	if err != nil {
		result.Success = false
		result.ErrorMessage = err.Error()
		result.ExitCode = 1
	} else {
		result.Success = true
		result.ExitCode = 0
	}

	return result, nil
}

// TFPlan executes opentofu plan
func (impl *TerraformStationImpl) TFPlan(ctx context.Context, input *TerraformStation.TFCommandInput) (*TerraformStation.TFPlanResult, error) {
	// Override command to ensure it's plan
	input.Command = "plan"

	// Execute plan command
	result, err := impl.TFCommand(ctx, input)
	if err != nil {
		return nil, err
	}

	// Parse plan output to determine if there are changes
	hasChanges := parsePlanOutput(result.Result)
	resourceCount := countResourcesInPlan(result.Result)

	planResult := &TerraformStation.TFPlanResult{
		PlanId:       TerraformStation.GenerateCommandID(),
		PlanOutput:   result.Result,
		HasChanges:   hasChanges,
		ResourceCount: int32(resourceCount),
		CreatedAt:    timestamppb.Now(),
		Status:       "completed",
	}

	if !result.Success {
		planResult.Status = "failed"
	}

	return planResult, nil
}

// TFApply executes opentofu apply
func (impl *TerraformStationImpl) TFApply(ctx context.Context, input *TerraformStation.TFCommandInput) (*TerraformStation.TFApplyResult, error) {
	// Override command to ensure it's apply
	input.Command = "apply"

	// Execute apply command
	result, err := impl.TFCommand(ctx, input)
	if err != nil {
		return nil, err
	}

	// Parse apply output to count resources
	resourcesAdded, resourcesChanged, resourcesDestroyed := parseApplyOutput(result.Result)

	applyResult := &TerraformStation.TFApplyResult{
		ApplyId:           TerraformStation.GenerateCommandID(),
		ApplyOutput:       result.Result,
		Success:           result.Success,
		ResourcesAdded:    int32(resourcesAdded),
		ResourcesChanged:  int32(resourcesChanged),
		ResourcesDestroyed: int32(resourcesDestroyed),
		ExecutedAt:        timestamppb.Now(),
	}

	return applyResult, nil
}

// TFInit executes opentofu init
func (impl *TerraformStationImpl) TFInit(ctx context.Context, input *TerraformStation.TFCommandInput) (*TerraformStation.TFCommandResult, error) {
	input.Command = "init"
	return impl.TFCommand(ctx, input)
}

// TFValidate executes opentofu validate
func (impl *TerraformStationImpl) TFValidate(ctx context.Context, input *TerraformStation.TFCommandInput) (*TerraformStation.TFCommandResult, error) {
	input.Command = "validate"
	return impl.TFCommand(ctx, input)
}

// TFState retrieves opentofu state information
func (impl *TerraformStationImpl) TFState(ctx context.Context, input *TerraformStation.TFCommandInput) (*TerraformStation.TFStateInfo, error) {
	// Use opentofu show to get state information
	input.Command = "show"
	
	result, err := impl.TFCommand(ctx, input)
	if err != nil {
		return nil, err
	}

	// Parse state output
	resourceCount := countResourcesInState(result.Result)
	terraformVersion := extractOpenTofuVersion(result.Result)

	stateInfo := &TerraformStation.TFStateInfo{
		StateId:          TerraformStation.GenerateCommandID(),
		StateFile:        input.StateFile,
		ResourceCount:    int32(resourceCount),
		LastUpdated:      timestamppb.Now(),
		TerraformVersion: terraformVersion,
	}

	return stateInfo, nil
}

// GetConfig returns the current configuration
func (impl *TerraformStationImpl) GetConfig() *TerraformStation.Config {
	return impl.cfg
}

// SetWorkingDirectory sets the working directory for OpenTofu operations
func (impl *TerraformStationImpl) SetWorkingDirectory(dir string) error {
	if err := impl.ValidateWorkingDirectory(dir); err != nil {
		return err
	}
	impl.workingDir = dir
	return nil
}

// ValidateWorkingDirectory validates the working directory
func (impl *TerraformStationImpl) ValidateWorkingDirectory(dir string) error {
	return impl.executor.ValidateWorkingDirectory(dir)
}

// Helper functions for parsing OpenTofu output

func parsePlanOutput(output string) bool {
	// Simple heuristic: if output contains "No changes" then no changes
	return !contains(output, "No changes")
}

func countResourcesInPlan(output string) int {
	// Simple counting of resource mentions in plan output
	// This is a basic implementation - in production you might want more sophisticated parsing
	count := 0
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "+ ") || strings.Contains(line, "- ") || strings.Contains(line, "~ ") {
			count++
		}
	}
	return count
}

func parseApplyOutput(output string) (added, changed, destroyed int) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Resources:") {
			// Parse the resources summary line
			// Example: "Resources: 2 added, 0 changed, 0 destroyed"
			// This is a simplified parser - in production you'd want more robust parsing
			if strings.Contains(line, "added") {
				added = 1 // Simplified for now
			}
			if strings.Contains(line, "changed") {
				changed = 1 // Simplified for now
			}
			if strings.Contains(line, "destroyed") {
				destroyed = 1 // Simplified for now
			}
		}
	}
	return
}

func countResourcesInState(output string) int {
	// Count resource blocks in state output
	count := 0
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "resource \"") {
			count++
		}
	}
	return count
}

func extractOpenTofuVersion(output string) string {
	// Extract OpenTofu version from output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "OpenTofu v") {
			// Extract version number
			parts := strings.Split(line, "OpenTofu v")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
		// Also check for "Terraform v" as OpenTofu might still show this
		if strings.Contains(line, "Terraform v") {
			parts := strings.Split(line, "Terraform v")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "unknown"
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
