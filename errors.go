package TerraformStation

import (
	"fmt"
	"strings"
)

// Error types for Terraform operations
type TerraformError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *TerraformError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common error codes
const (
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodeExecutionFailed  = "EXECUTION_FAILED"
	ErrCodeTimeout          = "TIMEOUT"
	ErrCodeInvalidState     = "INVALID_STATE"
	ErrCodeWorkingDirError  = "WORKING_DIR_ERROR"
	ErrCodeTerraformNotFound = "TERRAFORM_NOT_FOUND"
	ErrCodePermissionDenied = "PERMISSION_DENIED"
)

// Error constructors
func NewInvalidInputError(message string, details ...string) *TerraformError {
	return &TerraformError{
		Code:    ErrCodeInvalidInput,
		Message: message,
		Details: strings.Join(details, "; "),
	}
}

func NewExecutionFailedError(message string, details ...string) *TerraformError {
	return &TerraformError{
		Code:    ErrCodeExecutionFailed,
		Message: message,
		Details: strings.Join(details, "; "),
	}
}

func NewTimeoutError(message string, details ...string) *TerraformError {
	return &TerraformError{
		Code:    ErrCodeTimeout,
		Message: message,
		Details: strings.Join(details, "; "),
	}
}

func NewWorkingDirError(message string, details ...string) *TerraformError {
	return &TerraformError{
		Code:    ErrCodeWorkingDirError,
		Message: message,
		Details: strings.Join(details, "; "),
	}
}

func NewTerraformNotFoundError(message string, details ...string) *TerraformError {
	return &TerraformError{
		Code:    ErrCodeTerraformNotFound,
		Message: message,
		Details: strings.Join(details, "; "),
	}
}

func NewPermissionDeniedError(message string, details ...string) *TerraformError {
	return &TerraformError{
		Code:    ErrCodePermissionDenied,
		Message: message,
		Details: strings.Join(details, "; "),
	}
}
