# spemu - Spanner Emulator DML Inserter

A command-line tool for inserting DML (Data Manipulation Language) statements into Google Cloud Spanner Emulator.

## Features

- Parse and execute DML files (INSERT, UPDATE, DELETE statements)
- Support for SQL comments (`--` style)
- Dry run mode for validation
- Verbose output for debugging
- Integration with Spanner Emulator
- Comprehensive test suite with CI/CD

## Installation

### Install from source

```bash
go install github.com/nu0ma/spemu@latest
```

### Build from source

```bash
git clone https://github.com/nu0ma/spemu.git
cd spemu
make build
```

## Usage

### Basic Usage

```bash
spemu [options] <dml-file>
```

### Options

- `--project, -p`: Spanner project ID (required)
- `--instance, -i`: Spanner instance ID (required)  
- `--database, -d`: Spanner database ID (required)
- `--port, -P`: Spanner emulator port (default: 9010)
- `--dry-run`: Parse and validate DML without executing
- `--verbose`: Enable verbose output
- `--help`: Show help message

### Examples

```bash
# Execute DML file against emulator
spemu -p test-project -i test-instance -d test-database ./examples/seed.sql

# Dry run to validate SQL
spemu -p test-project -i test-instance -d test-database --dry-run ./examples/seed.sql

# Verbose output
spemu -p test-project -i test-instance -d test-database --verbose ./examples/seed.sql

# Use custom emulator port
spemu -p test-project -i test-instance -d test-database --port=9020 ./examples/seed.sql
```

## DML File Format

spemu supports SQL files with:
- `INSERT`, `UPDATE`, `DELETE` statements
- SQL comments using `--`
- Multiple statements separated by semicolons

Example DML file:
```sql
-- Insert users
INSERT INTO users (id, name, email) VALUES 
  (1, 'John Doe', 'john@example.com'),
  (2, 'Jane Smith', 'jane@example.com');

-- Update user
UPDATE users SET name = 'John Updated' WHERE id = 1;

-- Delete user
DELETE FROM users WHERE id = 2;
```

## Development

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Make (optional, for convenience commands)

### Setup Development Environment

```bash
# Clone repository
git clone https://github.com/nu0ma/spemu.git
cd spemu

# Install development dependencies
make dev-setup

# Start Spanner emulator
make emulator-setup
```

### Running Tests

```bash
# Run unit tests
make test-unit

# Run integration tests (requires emulator)
make test-integration

# Run all tests
make test-all

# Run tests with coverage
make test-coverage
```

### Development Workflow

```bash
# Format code, lint, and run unit tests
make dev-test

# Format code, lint, and run all tests
make dev-full

# Run demo with example data
make demo
```

### Emulator Management

```bash
# Start emulator and initialize database
make emulator-setup

# Start emulator only
make emulator-start

# Initialize database schema
make emulator-init

# Stop emulator
make emulator-stop

# Reset emulator (stop, start, initialize)
make emulator-reset

# View emulator logs
make emulator-logs
```

## Docker Usage

### Start Emulator

```bash
# Start emulator
docker-compose up -d spanner-emulator

# Initialize database
docker-compose --profile init up spanner-init
```

### Development Container

```bash
# Start development environment
docker-compose --profile dev up -d spemu-dev

# Run tests in container
docker-compose exec spemu-dev make test-unit
```

## CI/CD

The project includes GitHub Actions workflows for:
- Unit and integration testing
- Multi-platform builds (Linux, macOS, Windows)
- Code quality checks (formatting, linting, vetting)
- Test coverage reporting

## Project Structure

```
├── main.go              # Main application
├── pkg/                 # Library packages
│   ├── config/          # Configuration handling
│   ├── executor/        # Spanner execution logic
│   └── parser/          # DML parsing logic
├── test/                # Integration tests and test data
│   ├── schema.sql       # Test database schema
│   └── integration_test.go
├── examples/            # Example DML files
│   └── seed.sql
├── .github/workflows/   # CI/CD workflows
├── docker-compose.yml   # Docker development environment
├── Makefile            # Development commands
└── README.md
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make dev-full`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Google Cloud Spanner](https://cloud.google.com/spanner)
- [Spanner Emulator](https://cloud.google.com/spanner/docs/emulator)
- [Go Client Library for Spanner](https://pkg.go.dev/cloud.google.com/go/spanner)