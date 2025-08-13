package TerraformStation

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// OpenTofuExecutor handles the execution of OpenTofu commands
type OpenTofuExecutor struct {
	opentofuPath string
	timeout       time.Duration
}

// NewOpenTofuExecutor creates a new OpenTofu executor
func NewOpenTofuExecutor(opentofuPath string, timeout time.Duration) *OpenTofuExecutor {
	return &OpenTofuExecutor{
		opentofuPath: opentofuPath,
		timeout:       timeout,
	}
}

// Execute runs an OpenTofu command with the given arguments
func (e *OpenTofuExecutor) Execute(ctx context.Context, workingDir string, args ...string) (string, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// Validate working directory
	if err := e.ValidateWorkingDirectory(workingDir); err != nil {
		return "", err
	}

	// Check if opentofu binary exists
	if err := e.checkOpenTofuBinary(); err != nil {
		return "", err
	}

	// Prepare command
	cmd := exec.CommandContext(ctx, e.opentofuPath, args...)
	cmd.Dir = workingDir
	cmd.Env = os.Environ()

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("opentofu command failed: %w", err)
	}

	return string(output), nil
}

// ValidateWorkingDirectory checks if the working directory is valid
func (e *OpenTofuExecutor) ValidateWorkingDirectory(dir string) error {
	if dir == "" {
		return NewWorkingDirError("working directory cannot be empty")
	}

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return NewWorkingDirError("working directory does not exist", dir)
	}

	// Check if it's a directory
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return NewWorkingDirError("cannot access working directory", err.Error())
	}

	if !fileInfo.IsDir() {
		return NewWorkingDirError("path is not a directory", dir)
	}

	return nil
}

// checkOpenTofuBinary verifies that the opentofu binary exists and is executable
func (e *OpenTofuExecutor) checkOpenTofuBinary() error {
	if e.opentofuPath == "" {
		return NewTerraformNotFoundError("opentofu path is not set")
	}

	// Check if binary exists
	if _, err := os.Stat(e.opentofuPath); os.IsNotExist(err) {
		return NewTerraformNotFoundError("opentofu binary not found", e.opentofuPath)
	}

	// Check if it's executable
	if _, err := exec.LookPath(e.opentofuPath); err != nil {
		return NewTerraformNotFoundError("opentofu binary is not executable", e.opentofuPath)
	}

	return nil
}

// BuildOpenTofuArgs constructs the arguments for an OpenTofu command
func BuildOpenTofuArgs(command string, input *TFCommandInput) []string {
	args := []string{command}

	// Add working directory if specified
	if input.WorkingDirectory != "" {
		args = append(args, "-chdir="+input.WorkingDirectory)
	}

	// Add variables
	for key, value := range input.Variables {
		args = append(args, "-var", fmt.Sprintf("%s=%s", key, value))
	}

	// Add additional arguments
	args = append(args, input.Arguments...)

	// Add plan file if specified
	if input.PlanFile != "" {
		args = append(args, input.PlanFile)
	}

	// Add state file if specified
	if input.StateFile != "" {
		args = append(args, "-state="+input.StateFile)
	}

	return args
}

// ValidateTFCommandInput validates the input for OpenTofu commands
func ValidateTFCommandInput(input *TFCommandInput) error {
	if input == nil {
		return NewInvalidInputError("input cannot be nil")
	}

	if input.Command == "" {
		return NewInvalidInputError("command cannot be empty")
	}

	// Validate command
	validCommands := map[string]bool{
		"init":     true,
		"plan":     true,
		"apply":    true,
		"destroy":  true,
		"validate": true,
		"state":    true,
		"output":   true,
		"show":     true,
		"version":  true,
	}

	if !validCommands[input.Command] {
		return NewInvalidInputError("invalid opentofu command", input.Command)
	}

	return nil
}

// SanitizeWorkingDirectory ensures the working directory path is safe
func SanitizeWorkingDirectory(dir string) string {
	// Clean the path
	cleanPath := filepath.Clean(dir)
	
	// Ensure it's not an absolute path that could be dangerous
	if filepath.IsAbs(cleanPath) {
		// For now, just return the cleaned path
		// In production, you might want to restrict to specific directories
		return cleanPath
	}
	
	return cleanPath
}

// GenerateCommandID creates a unique identifier for a command execution
func GenerateCommandID() string {
	return fmt.Sprintf("tofu_%d", time.Now().UnixNano())
}
