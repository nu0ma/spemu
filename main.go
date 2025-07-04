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
		dryRun    = flag.Bool("dry-run", false, "Parse and validate DML without executing")
		verbose   = flag.Bool("verbose", false, "Enable verbose output")
		help      = flag.Bool("help", false, "Show help message")
		version   = flag.Bool("version", false, "Show version information")
		project   = flag.String("project", "", "Spanner project ID (required)")
		projectS  = flag.String("p", "", "Spanner project ID (short form)")
		instance  = flag.String("instance", "", "Spanner instance ID (required)")
		instanceS = flag.String("i", "", "Spanner instance ID (short form)")
		database  = flag.String("database", "", "Spanner database ID (required)")
		databaseS = flag.String("d", "", "Spanner database ID (short form)")
		port      = flag.String("port", "9010", "Spanner emulator port (default: 9010)")
		portS     = flag.String("P", "9010", "Spanner emulator port (short form)")
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

	// Use short form if long form is empty
	if *project == "" && *projectS != "" {
		*project = *projectS
	}
	if *instance == "" && *instanceS != "" {
		*instance = *instanceS
	}
	if *database == "" && *databaseS != "" {
		*database = *databaseS
	}
	if *port == "9010" && *portS != "9010" {
		*port = *portS
	}

	// Validate required flags
	if *project == "" {
		fmt.Fprintf(os.Stderr, "Error: --project (or -p) is required\n")
		os.Exit(1)
	}
	if *instance == "" {
		fmt.Fprintf(os.Stderr, "Error: --instance (or -i) is required\n")
		os.Exit(1)
	}
	if *database == "" {
		fmt.Fprintf(os.Stderr, "Error: --database (or -d) is required\n")
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
			fmt.Printf("Statement %d: %s\n", i+1, stmt[:min(50, len(stmt))]+"...")
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
  --project, -p    Spanner project ID (required)
  --instance, -i   Spanner instance ID (required)
  --database, -d   Spanner database ID (required)
  --port, -P       Spanner emulator port (default: 9010)
  --dry-run        Parse and validate DML without executing
  --verbose        Enable verbose output
  --version        Show version information
  --help           Show this help message

Examples:
  spemu -p test-project -i test-instance -d test-database ./seed.sql
  spemu --project=my-proj --instance=my-inst --database=my-db --dry-run ./test.sql
  spemu -p test -i test -d test --port=9020 ./users.sql

`)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
