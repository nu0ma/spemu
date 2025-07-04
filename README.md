# spemu - Spanner Emulator DML Inserter

A simple CLI tool to execute DML (INSERT, UPDATE, DELETE) statements against Google Cloud Spanner Emulator.

## Installation

```bash
go build -o spemu cmd/spemu/main.go
```

## Usage

### Basic Usage
```bash
spemu -p test-project -i test-instance -d test-database ./seed.sql
```

### Options
```bash
spemu -p my-proj -i my-inst -d my-db --dry-run ./seed.sql      # Parse and validate without executing
spemu -p my-proj -i my-inst -d my-db --verbose ./seed.sql      # Enable verbose output
spemu --help                                                   # Show help message
```

## Required Flags

spemu requires the following flags to connect to Spanner emulator:

- `--project, -p`: Spanner project ID (required)
- `--instance, -i`: Spanner instance ID (required) 
- `--database, -d`: Spanner database ID (required)
- `--host`: Spanner emulator host (optional, default: localhost:9010)

## DML File Format

spemu accepts SQL files containing DML statements separated by semicolons:

```sql
-- Comments are supported
INSERT INTO users (id, name, email, created_at) VALUES
  (1, 'John Doe', 'john.doe@example.com', '2024-01-01T00:00:00Z'),
  (2, 'Jane Smith', 'jane.smith@example.com', '2024-01-02T00:00:00Z');

INSERT INTO posts (id, user_id, title, content) VALUES
  (1, 1, 'First Post', 'Hello World! Testing with Spanner emulator.'),
  (2, 2, 'About Spanner', 'Cloud Spanner is an amazing database service.');
```

### Supported Statements
- `INSERT` statements
- `UPDATE` statements  
- `DELETE` statements

## Example Workflow

1. Start Spanner Emulator:
```bash
gcloud spanner emulators start
```

2. Create instance and database (if needed):
```bash
gcloud spanner instances create test-instance --config=emulator-config --description="Test Instance"
gcloud spanner databases create test-database --instance=test-instance
```

3. Execute seed data:
```bash
spemu -p test-project -i test-instance -d test-database ./examples/seed.sql
```

## Features

- **Transaction Safety**: All statements are executed within a single transaction
- **Error Handling**: Detailed error messages with statement context
- **Dry Run**: Validate DML files without executing
- **Verbose Mode**: Progress reporting during execution
- **Comment Support**: SQL-style comments are stripped from input
- **Emulator Optimized**: Designed specifically for Spanner Emulator

## License

MIT