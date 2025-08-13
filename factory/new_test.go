package factory

import (
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

	// Test successful creation
	service := New(db, cfg)
	assert.NotNil(t, service)
	assert.Implements(t, (*TerraformStation.TerraformStationService)(nil), service)
}

func TestNewWithNilDB(t *testing.T) {
	cfg := TerraformStation.DefaultConfig()
	
	// This should panic with nil database
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil database")
		}
	}()
	
	New(nil, cfg)
}

func TestNewWithNilConfig(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	
	// This should panic with nil config
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil config")
		}
	}()
	
	New(db, nil)
}
