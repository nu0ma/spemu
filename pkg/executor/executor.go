package executor

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"github.com/nu0ma/spemu/pkg/config"
)

type Executor struct {
	client *spanner.Client
}

func New(cfg *config.Config) (*Executor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if cfg.EmulatorHost != "" {
		os.Setenv("SPANNER_EMULATOR_HOST", cfg.EmulatorHost)
	}

	client, err := spanner.NewClient(ctx, cfg.DatabasePath())
	if err != nil {
		return nil, fmt.Errorf("failed to create Spanner client: %w", err)
	}

	return &Executor{client: client}, nil
}

func (e *Executor) Close() {
	if e.client != nil {
		e.client.Close()
	}
}

func (e *Executor) ExecuteStatements(statements []string, verbose bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := e.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		for i, stmt := range statements {
			if verbose {
				limit := 100
				if len(stmt) < limit {
					limit = len(stmt)
				}
				fmt.Printf("Executing statement %d/%d: %s\n", i+1, len(statements), stmt[:limit]+"...")
			}

			_, err := txn.Update(ctx, spanner.Statement{SQL: stmt})
			if err != nil {
				return fmt.Errorf("failed to execute statement %d: %w\nStatement: %s", i+1, err, stmt)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

// InitializeSchema creates instance and database with the given schema
func InitializeSchema(cfg *config.Config, schemaFile string, verbose bool) error {
	if cfg.EmulatorHost != "" {
		os.Setenv("SPANNER_EMULATOR_HOST", cfg.EmulatorHost)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create instance admin client
	instanceAdminClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create instance admin client: %w", err)
	}
	defer instanceAdminClient.Close()

	// Create database admin client
	databaseAdminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create database admin client: %w", err)
	}
	defer databaseAdminClient.Close()

	// Create instance first (for emulator)
	projectPath := fmt.Sprintf("projects/%s", cfg.ProjectID)
	instancePath := fmt.Sprintf("projects/%s/instances/%s", cfg.ProjectID, cfg.InstanceID)

	// Check if instance exists
	_, err = instanceAdminClient.GetInstance(ctx, &instancepb.GetInstanceRequest{
		Name: instancePath,
	})

	if err != nil {
		// Instance doesn't exist, create it
		if verbose {
			fmt.Printf("Creating instance: %s\n", instancePath)
		}

		createInstanceOp, err := instanceAdminClient.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
			Parent:     projectPath,
			InstanceId: cfg.InstanceID,
			Instance: &instancepb.Instance{
				Name:        instancePath,
				DisplayName: cfg.InstanceID,
				NodeCount:   1,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create instance: %w", err)
		}

		// Wait for instance creation to complete
		_, err = createInstanceOp.Wait(ctx)
		if err != nil {
			return fmt.Errorf("instance creation failed: %w", err)
		}

		if verbose {
			fmt.Printf("Instance created successfully: %s\n", cfg.InstanceID)
		}
	} else {
		if verbose {
			fmt.Printf("Instance already exists: %s\n", cfg.InstanceID)
		}
	}

	// Check if database exists
	databasePath := cfg.DatabasePath()
	_, err = databaseAdminClient.GetDatabase(ctx, &databasepb.GetDatabaseRequest{
		Name: databasePath,
	})

	if err != nil {
		// Database doesn't exist, create it
		if verbose {
			fmt.Printf("Creating database: %s\n", databasePath)
		}

		// Read schema file
		schemaContent, err := os.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}

		// Parse DDL statements
		ddlStatements := parseDDLStatements(string(schemaContent))

		if verbose {
			fmt.Printf("Found %d DDL statements\n", len(ddlStatements))
		}

		// Create database with schema
		createOp, err := databaseAdminClient.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
			Parent:          instancePath,
			CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", cfg.DatabaseID),
			ExtraStatements: ddlStatements,
		})
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}

		// Wait for database creation to complete
		_, err = createOp.Wait(ctx)
		if err != nil {
			return fmt.Errorf("database creation failed: %w", err)
		}

		if verbose {
			fmt.Printf("Database created successfully: %s\n", cfg.DatabaseID)
		}
	} else {
		if verbose {
			fmt.Printf("Database already exists: %s\n", cfg.DatabaseID)
		}
	}

	return nil
}

// parseDDLStatements parses DDL statements from schema content
func parseDDLStatements(content string) []string {
	// Remove comment lines
	lines := strings.Split(content, "\n")
	var cleanLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comment lines
		if line != "" && !strings.HasPrefix(line, "--") {
			cleanLines = append(cleanLines, line)
		}
	}

	// Join back and split by semicolon
	cleanContent := strings.Join(cleanLines, " ")
	statements := strings.Split(cleanContent, ";")
	var ddlStatements []string

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			ddlStatements = append(ddlStatements, stmt)
		}
	}

	return ddlStatements
}
