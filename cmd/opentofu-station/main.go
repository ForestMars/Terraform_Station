package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ForestMars/TerraformStation"
	"github.com/ForestMars/TerraformStation/factory"
)

func main() {
	// Parse command line flags
	opentofuPath := flag.String("opentofu", "tofu", "Path to opentofu binary")
	workingDir := flag.String("workdir", "./tofu", "Working directory for OpenTofu operations")
	port := flag.String("port", "8080", "Port to listen on")
	host := flag.String("host", "localhost", "Host to bind to")
	dbDriver := flag.String("db-driver", "sqlite", "Database driver (sqlite or postgres)")
	flag.Parse()

	// Create default configuration
	cfg := TerraformStation.DefaultConfig()
	
	// Override with command line flags
	if *opentofuPath != "" {
		cfg.OpenTofuPath = *opentofuPath
	}
	if *workingDir != "" {
		cfg.WorkingDirectory = *workingDir
	}
	if *port != "" {
		cfg.Port = *port
	}
	if *host != "" {
		cfg.Host = *host
	}
	if *dbDriver != "" {
		cfg.Database.Driver = *dbDriver
	}

	// Initialize database
	dbManager, err := TerraformStation.NewDatabaseManager(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbManager.Close()

	// Get GORM database instance
	db := dbManager.GetDB()

	// Create service instance
	service := factory.New(db, cfg)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, shutting down gracefully...")
		cancel()
	}()

	// Start the service
	log.Printf("Starting OpenTofu Station on %s:%s", cfg.Host, cfg.Port)
	log.Printf("OpenTofu binary: %s", cfg.OpenTofuPath)
	log.Printf("Working directory: %s", cfg.WorkingDirectory)
	log.Printf("Database driver: %s", cfg.Database.Driver)

	// For now, just run a simple test command
	if err := runTestCommand(ctx, service); err != nil {
		log.Printf("Test command failed: %v", err)
	}

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("OpenTofu Station stopped")
}

func runTestCommand(ctx context.Context, service TerraformStation.TerraformStationService) error {
	log.Println("Running test command...")

	// Create a test input
	input := &TerraformStation.TFCommandInput{
		Command:        "version",
		WorkingDirectory: "./tofu",
	}

	// Execute the command
	result, err := service.TFCommand(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to execute opentofu version: %w", err)
	}

	log.Printf("Command executed successfully: %s", result.Result)
	return nil
}
