package internal

import (
	"context"
	"testing"

	"github.com/ForestMars/TerraformStation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	// Create a test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create test configuration
	cfg := TerraformStation.DefaultConfig()
	cfg.WorkingDirectory = "/tmp" // Use a directory that exists

	// Test successful creation
	impl, err := New(db, cfg)
	require.NoError(t, err)
	assert.NotNil(t, impl)
	assert.Equal(t, cfg, impl.cfg)
	assert.Equal(t, db, impl.db)
}

func TestNewWithNilConfig(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Test with nil config
	impl, err := New(db, nil)
	assert.Error(t, err)
	assert.Nil(t, impl)
}

func TestNewWithNilDB(t *testing.T) {
	cfg := TerraformStation.DefaultConfig()

	// Test with nil database
	impl, err := New(nil, cfg)
	assert.Error(t, err)
	assert.Nil(t, impl)
}

func TestTFCommand(t *testing.T) {
	// Create a test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create test configuration
	cfg := TerraformStation.DefaultConfig()
	cfg.OpenTofuPath = "echo" // Use echo for testing
	cfg.WorkingDirectory = "/tmp"

	// Create implementation
	impl, err := New(db, cfg)
	require.NoError(t, err)

	// Test valid command
	input := &TerraformStation.TFCommandInput{
		Command: "version",
	}

	ctx := context.Background()
	result, err := impl.TFCommand(ctx, input)
	
	// Since we're using echo, this should fail, but we can test the structure
	if err != nil {
		// Expected error for invalid terraform command
		assert.Contains(t, err.Error(), "terraform")
	} else {
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.CommandId)
		assert.NotNil(t, result.ExecutedAt)
	}
}

func TestValidateWorkingDirectory(t *testing.T) {
	// Create a test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create test configuration
	cfg := TerraformStation.DefaultConfig()

	// Create implementation
	impl, err := New(db, cfg)
	require.NoError(t, err)

	// Test valid directory
	err = impl.ValidateWorkingDirectory("/tmp")
	assert.NoError(t, err)

	// Test invalid directory
	err = impl.ValidateWorkingDirectory("/nonexistent/directory")
	assert.Error(t, err)
}

func TestGetConfig(t *testing.T) {
	// Create a test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create test configuration
	cfg := TerraformStation.DefaultConfig()

	// Create implementation
	impl, err := New(db, cfg)
	require.NoError(t, err)

	// Test getting config
	retrievedCfg := impl.GetConfig()
	assert.Equal(t, cfg, retrievedCfg)
}
