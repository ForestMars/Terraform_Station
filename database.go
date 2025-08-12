package TerraformStation

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseManager handles database operations
type DatabaseManager struct {
	db *gorm.DB
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager(cfg *Config) (*DatabaseManager, error) {
	var db *gorm.DB
	var err error

	switch cfg.Database.Driver {
	case "postgres":
		db, err = connectPostgres(cfg.Database)
	case "sqlite":
		db, err = connectSQLite(cfg.Database)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure GORM
	db.Logger = db.Logger.LogMode(logger.Info)

	// Auto migrate models
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	return &DatabaseManager{db: db}, nil
}

// GetDB returns the underlying GORM database instance
func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.db
}

// Close closes the database connection
func (dm *DatabaseManager) Close() error {
	sqlDB, err := dm.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// connectPostgres establishes a connection to PostgreSQL
func connectPostgres(cfg DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// connectSQLite establishes a connection to SQLite
func connectSQLite(cfg DatabaseConfig) (*gorm.DB, error) {
	// For SQLite, we'll use the database name as the file path
	// If no database name is specified, use a default
	dbPath := cfg.Database
	if dbPath == "" {
		dbPath = "terraform_station.db"
	}

	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}

// autoMigrate automatically migrates the database schema
func autoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Migrate all models
	err := db.AutoMigrate(
		&TerraformOperation{},
		&TerraformPlan{},
		&TerraformApply{},
		&TerraformState{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// CreateOperation creates a new Terraform operation record
func (dm *DatabaseManager) CreateOperation(operation *TerraformOperation) error {
	return dm.db.Create(operation).Error
}

// UpdateOperation updates an existing Terraform operation
func (dm *DatabaseManager) UpdateOperation(operation *TerraformOperation) error {
	return dm.db.Save(operation).Error
}

// GetOperationByID retrieves an operation by its ID
func (dm *DatabaseManager) GetOperationByID(id uint) (*TerraformOperation, error) {
	var operation TerraformOperation
	err := dm.db.First(&operation, id).Error
	if err != nil {
		return nil, err
	}
	return &operation, nil
}

// GetOperationByCommandID retrieves an operation by its command ID
func (dm *DatabaseManager) GetOperationByCommandID(commandID string) (*TerraformOperation, error) {
	var operation TerraformOperation
	err := dm.db.Where("command_id = ?", commandID).First(&operation).Error
	if err != nil {
		return nil, err
	}
	return &operation, nil
}

// ListOperations retrieves a list of operations with optional filtering
func (dm *DatabaseManager) ListOperations(limit, offset int, status string) ([]TerraformOperation, error) {
	var operations []TerraformOperation
	query := dm.db

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&operations).Error
	return operations, err
}

// CreatePlan creates a new Terraform plan record
func (dm *DatabaseManager) CreatePlan(plan *TerraformPlan) error {
	return dm.db.Create(plan).Error
}

// CreateApply creates a new Terraform apply record
func (dm *DatabaseManager) CreateApply(apply *TerraformApply) error {
	return dm.db.Create(apply).Error
}

// CreateState creates a new Terraform state record
func (dm *DatabaseManager) CreateState(state *TerraformState) error {
	return dm.db.Create(state).Error
}
