# spemu - Spanner Emulator DML Inserter

A command-line tool for inserting DML (Data Manipulation Language) statements into Google Cloud Spanner Emulator.

## Features

- No dependency on `gcloud` CLI
- Parse and execute DML files (INSERT, UPDATE, DELETE statements)
- Support for SQL comments (`--` style)
- Dry run mode for validation
- Verbose output for debugging
- Integration with Spanner Emulator
- Comprehensive test suite with CI/CD

## Installation

### Install from Go (Recommended)

Install the latest release:
```bash
go install github.com/nu0ma/spemu@latest
```

Install a specific version:
```bash
go install github.com/nu0ma/spemu@v0.4.1
```

### Download Pre-built Binaries

Download the latest release from [GitHub Releases](https://github.com/nu0ma/spemu/releases):

#### Linux (x86_64)
```bash
curl -L https://github.com/nu0ma/spemu/releases/latest/download/spemu_Linux_x86_64.tar.gz | tar xz
sudo mv spemu /usr/local/bin/
```

#### macOS (Intel)
```bash
curl -L https://github.com/nu0ma/spemu/releases/latest/download/spemu_Darwin_x86_64.tar.gz | tar xz
sudo mv spemu /usr/local/bin/
```

#### macOS (Apple Silicon)
```bash
curl -L https://github.com/nu0ma/spemu/releases/latest/download/spemu_Darwin_arm64.tar.gz | tar xz
sudo mv spemu /usr/local/bin/
```

#### Windows
Download `spemu_Windows_x86_64.zip` from the [releases page](https://github.com/nu0ma/spemu/releases) and extract the executable.

### Using Homebrew (macOS/Linux)

**Coming Soon** - Homebrew tap will be available in a future release:
```bash
brew install nu0ma/tap/spemu
```

### Build from Source

```bash
git clone https://github.com/nu0ma/spemu.git
cd spemu
make build
```

### Verification

Verify the installation:
```bash
spemu --version
```

All release binaries include checksums for verification:
```bash
# Download checksums
curl -L https://github.com/nu0ma/spemu/releases/latest/download/checksums.txt

# Verify binary (example for Linux)
sha256sum -c checksums.txt --ignore-missing
```

## Usage

### Basic Usage

```bash
spemu [options] <dml-file>
```

### Options

- `--project`: Spanner project ID (required)
- `--instance`: Spanner instance ID (required)  
- `--database`: Spanner database ID (required)
- `--port`: Spanner emulator port (default: 9010)
- `--dry-run`: Parse and validate DML without executing
- `--verbose`: Enable verbose output
- `--help`: Show help message

### Examples

```bash
# Execute DML file against emulator
spemu --project=test-project --instance=test-instance --database=test-database ./examples/seed.sql

# Dry run to validate SQL
spemu --project=test-project --instance=test-instance --database=test-database --dry-run ./examples/seed.sql
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

**Note:** spemu does not require the `gcloud` CLI tool. It connects directly to the Spanner Emulator using Go client libraries.

### Setup Development Environment

```bash
# Clone repository
git clone https://github.com/nu0ma/spemu.git
cd spemu

# Start Spanner emulator
make emulator-setup
```

### Running Tests

```bash
# Run unit tests
make test-unit

# Run integration tests (requires emulator)
make test-integration
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

We welcome contributions! This project uses automated versioning and releases.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests (`make test-unit`)
5. Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages
6. Push to the branch
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Google Cloud Spanner](https://cloud.google.com/spanner)
- [Spanner Emulator](https://cloud.google.com/spanner/docs/emulator)
- [Go Client Library for Spanner](https://pkg.go.dev/cloud.google.com/go/spanner)