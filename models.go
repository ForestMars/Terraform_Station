package TerraformStation

import (
	"time"

	"gorm.io/gorm"
)

// TerraformOperation represents a Terraform operation in the database
type TerraformOperation struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CommandID     string         `gorm:"uniqueIndex;not null" json:"command_id"`
	Command       string         `gorm:"not null" json:"command"`
	WorkingDir    string         `gorm:"not null" json:"working_dir"`
	Arguments     string         `gorm:"type:text" json:"arguments"`
	Variables     string         `gorm:"type:text" json:"variables"`
	Status        string         `gorm:"not null;default:'pending'" json:"status"`
	ExitCode      int            `gorm:"default:0" json:"exit_code"`
	Output        string         `gorm:"type:text" json:"output"`
	ErrorMessage  string         `gorm:"type:text" json:"error_message"`
	StartedAt     time.Time      `gorm:"not null" json:"started_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
	Duration      time.Duration  `json:"duration"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TerraformPlan represents a Terraform plan operation
type TerraformPlan struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	PlanID        string         `gorm:"uniqueIndex;not null" json:"plan_id"`
	OperationID   uint           `gorm:"not null" json:"operation_id"`
	Operation     TerraformOperation `gorm:"foreignKey:OperationID" json:"operation"`
	HasChanges    bool           `gorm:"not null" json:"has_changes"`
	ResourceCount int            `gorm:"default:0" json:"resource_count"`
	PlanOutput    string         `gorm:"type:text" json:"plan_output"`
	Status        string         `gorm:"not null;default:'pending'" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TerraformApply represents a Terraform apply operation
type TerraformApply struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	ApplyID           string         `gorm:"uniqueIndex;not null" json:"apply_id"`
	OperationID       uint           `gorm:"not null" json:"operation_id"`
	Operation         TerraformOperation `gorm:"foreignKey:OperationID" json:"operation"`
	PlanID            string         `json:"plan_id"`
	Success           bool           `gorm:"not null" json:"success"`
	ResourcesAdded    int            `gorm:"default:0" json:"resources_added"`
	ResourcesChanged  int            `gorm:"default:0" json:"resources_changed"`
	ResourcesDestroyed int           `gorm:"default:0" json:"resources_destroyed"`
	ApplyOutput       string         `gorm:"type:text" json:"apply_output"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TerraformState represents Terraform state information
type TerraformState struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	StateID           string         `gorm:"uniqueIndex;not null" json:"state_id"`
	StateFile         string         `gorm:"not null" json:"state_file"`
	WorkingDir        string         `gorm:"not null" json:"working_dir"`
	ResourceCount     int            `gorm:"default:0" json:"resource_count"`
	TerraformVersion  string         `json:"terraform_version"`
	LastUpdated       time.Time      `gorm:"not null" json:"last_updated"`
	StateData         string         `gorm:"type:text" json:"state_data"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for TerraformOperation
func (TerraformOperation) TableName() string {
	return "terraform_operations"
}

// TableName specifies the table name for TerraformPlan
func (TerraformPlan) TableName() string {
	return "terraform_plans"
}

// TableName specifies the table name for TerraformApply
func (TerraformApply) TableName() string {
	return "terraform_applies"
}

// TableName specifies the table name for TerraformState
func (TerraformState) TableName() string {
	return "terraform_states"
}
