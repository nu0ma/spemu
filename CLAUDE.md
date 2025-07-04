# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**spemu** is a Go command-line tool that executes DML (Data Manipulation Language) statements against Google Cloud Spanner Emulator. It provides atomic transaction execution for INSERT, UPDATE, and DELETE statements with comprehensive validation and error handling.

## Essential Development Commands

### Build and Development
```bash
make build                 # Build the binary
make dev-setup            # Install development dependencies (staticcheck, goimports)
make dev-test             # Format, lint, and run unit tests
```

### Testing
```bash
make test-unit            # Run unit tests with race detection and coverage
make test-integration     # Run integration tests (requires emulator)
make test-all            # Run all tests (unit + integration)
```

### Code Quality
```bash
make fmt                  # Format code with gofmt + goimports
make lint                 # Run staticcheck and go vet
make vet                  # Run go vet only
```

### Emulator Management
```bash
make emulator-setup       # Start emulator and initialize test database
make emulator-start       # Start Spanner emulator on localhost:9010
make emulator-stop        # Stop emulator
```

## Architecture Overview

### Core Package Structure
- **`pkg/config/`**: Configuration management with database path generation
- **`pkg/parser/`**: SQL parsing, comment removal, and DML validation 
- **`pkg/executor/`**: Spanner client lifecycle and transaction execution
- **`main.go`**: CLI entry point with flag-based configuration

### Key Architectural Patterns

**Transaction-Based Execution**: All DML statements execute within a single Spanner transaction. If any statement fails, the entire transaction rolls back atomically.

**DML-Only Validation**: The parser strictly validates that only INSERT, UPDATE, and DELETE statements are allowed. DDL and SELECT statements are rejected.

**Emulator Integration**: Uses `SPANNER_EMULATOR_HOST=localhost:9010` environment variable to connect to the local emulator instead of production Spanner.

**Modular Design**: Clean separation of concerns with dedicated packages for configuration, parsing, and execution logic.

## Development Environment Setup

### Prerequisites
- Go 1.24 or later
- Docker and Docker Compose
- Make

### Quick Start
```bash
# Setup development environment
make dev-setup

# Start emulator and initialize database
make emulator-setup

# Run tests to verify setup
make test-all

# Build and test CLI
make build
./spemu --help
```

## Testing Strategy

### Unit Tests
- Located alongside source code in `pkg/*/`
- Run with race detection enabled
- Coverage reporting via `coverage.out`

### Integration Tests
- Located in `test/integration_test.go`
- Requires running Spanner emulator
- Tests real database operations with comprehensive scenarios
- Automatic cleanup between tests

### Test Schema
- Defined in `test/schema.sql`
- Relational structure: users → posts → comments
- Uses Spanner-specific features (foreign keys, commit timestamps)

## CLI Usage Patterns

```bash
# Basic execution
spemu --project=test-project --instance=test-instance --database=test-database file.sql

# Dry run (validation only)
spemu --project=test --instance=test --database=test --dry-run file.sql

# Verbose output for debugging
spemu --verbose --project=test --instance=test --database=test file.sql
```

## Important Configuration Details

### Environment Variables
- `SPANNER_EMULATOR_HOST`: Set to `localhost:9010` for emulator usage
- `CI`: Used in tests to skip certain setup steps

### Version Management
- Version stored in `version.txt` (currently 0.3.0)
- Used by Makefile for build metadata
- Managed via release-please automation

### CI/CD Pipeline
- GitHub Actions workflows in `.github/workflows/`
- Runs code quality checks, unit tests, and integration tests
- Uses Docker for Spanner emulator in CI environment

## DML File Format

### Supported Statements
- INSERT, UPDATE, DELETE only
- Multiple statements separated by semicolons
- SQL comments with `--` style are supported and stripped during parsing

### Example Structure
```sql
-- Insert users
INSERT INTO users (id, name, email) VALUES (1, 'John', 'john@example.com');

-- Update user
UPDATE users SET name = 'Jane' WHERE id = 1;

-- Clean up
DELETE FROM users WHERE id = 1;
```

## Key Dependencies

- **`cloud.google.com/go/spanner`** (v1.83.0): Core Spanner client
- **`google.golang.org/api`** (v0.237.0): Google Cloud APIs
- **Development tools**: staticcheck, goimports for code quality