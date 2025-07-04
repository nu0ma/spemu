package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nu0ma/spemu/pkg/config"
	"github.com/nu0ma/spemu/pkg/executor"
	"github.com/nu0ma/spemu/pkg/parser"
)

// Version is set during build time via ldflags
var Version = "dev"

func main() {
	var (
		dryRun   = flag.Bool("dry-run", false, "Parse and validate DML without executing")
		verbose  = flag.Bool("verbose", false, "Enable verbose output")
		help     = flag.Bool("help", false, "Show help message")
		version  = flag.Bool("version", false, "Show version information")
		project  = flag.String("project", "", "Spanner project ID (required)")
		instance = flag.String("instance", "", "Spanner instance ID (required)")
		database = flag.String("database", "", "Spanner database ID (required)")
		port     = flag.String("port", "9010", "Spanner emulator port (default: 9010)")
	)
	flag.Parse()

	if *version {
		fmt.Printf("spemu version %s\n", Version)
		return
	}

	if *help {
		showHelp()
		return
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: spemu [options] <dml-file>\n")
		fmt.Fprintf(os.Stderr, "Run 'spemu --help' for more information.\n")
		os.Exit(1)
	}

	dmlFile := args[0]

	// Validate required flags
	if *project == "" {
		fmt.Fprintf(os.Stderr, "Error: --project is required\n")
		os.Exit(1)
	}
	if *instance == "" {
		fmt.Fprintf(os.Stderr, "Error: --instance is required\n")
		os.Exit(1)
	}
	if *database == "" {
		fmt.Fprintf(os.Stderr, "Error: --database is required\n")
		os.Exit(1)
	}

	emulatorHost := fmt.Sprintf("localhost:%s", *port)
	cfg := &config.Config{
		ProjectID:    *project,
		InstanceID:   *instance,
		DatabaseID:   *database,
		EmulatorHost: emulatorHost,
	}

	if *verbose {
		fmt.Printf("Configuration: %+v\n", cfg)
		fmt.Printf("DML file: %s\n", dmlFile)
	}

	statements, err := parser.ParseDMLFile(dmlFile)
	if err != nil {
		log.Fatalf("Failed to parse DML file: %v", err)
	}

	if *verbose {
		fmt.Printf("Parsed %d DML statements\n", len(statements))
	}

	if *dryRun {
		fmt.Printf("Dry run: %d statements would be executed\n", len(statements))
		for i, stmt := range statements {
			limit := 50
			if len(stmt) < limit {
				limit = len(stmt)
			}
			fmt.Printf("Statement %d: %s\n", i+1, stmt[:limit]+"...")
		}
		return
	}

	exec, err := executor.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create executor: %v", err)
	}
	defer exec.Close()

	err = exec.ExecuteStatements(statements, *verbose)
	if err != nil {
		log.Fatalf("Failed to execute statements: %v", err)
	}

	fmt.Printf("Successfully executed %d statements\n", len(statements))
}

func showHelp() {
	fmt.Printf(`spemu - Spanner Emulator DML Inserter

Usage:
  spemu [options] <dml-file>

Options:
  --project        Spanner project ID (required)
  --instance       Spanner instance ID (required)
  --database       Spanner database ID (required)
  --port           Spanner emulator port (default: 9010)
  --dry-run        Parse and validate DML without executing
  --verbose        Enable verbose output
  --version        Show version information
  --help           Show this help message

Examples:
  spemu --project=test-project --instance=test-instance --database=test-database ./seed.sql
  spemu --project=my-proj --instance=my-inst --database=my-db --dry-run ./test.sql
  spemu --project=test --instance=test --database=test --port=9020 ./users.sql

`)
}

