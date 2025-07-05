#!/bin/bash
set -e

PROJECT_ID="test-project"
INSTANCE_ID="test-instance"
DATABASE_ID="test-database"

# Function to check dependencies
check_dependencies() {
    echo "Checking dependencies..."
    
    if ! command -v docker &> /dev/null; then
        echo "Error: Docker is required but not installed"
        exit 1
    fi
    
    if ! command -v nc &> /dev/null; then
        echo "Error: nc (netcat) is required but not installed"
        echo "Please install netcat: apt-get install netcat (Debian/Ubuntu) or brew install netcat (macOS)"
        exit 1
    fi
    
    echo "All dependencies found ✓"
}

# Function to cleanup existing containers
cleanup_existing() {
    if docker ps -a --format 'table {{.Names}}' | grep -q spanner-emulator; then
        echo "Stopping and removing existing spanner-emulator container..."
        docker stop spanner-emulator >/dev/null 2>&1 || true
        docker rm spanner-emulator >/dev/null 2>&1 || true
    fi
}

# Function to wait for emulator to be ready
waitForEmulator() {
    local port=$1
    local max_attempts=30
    local attempt=1
    
    echo "Waiting for emulator to be ready on port $port..."
    
    while [ $attempt -le $max_attempts ]; do
        # Check if port is listening using nc (netcat)
        # -z: zero I/O mode (just check connectivity)
        # -w1: timeout of 1 second
        if nc -z -w1 localhost "$port" 2>/dev/null; then
            echo "Emulator is ready on port $port ✓"
            return 0
        fi
        
        # Show progress every 5 attempts
        if [ $((attempt % 5)) -eq 0 ]; then
            echo "Still waiting... (attempt $attempt/$max_attempts)"
        fi
        
        # Wait before retry
        sleep 1
        attempt=$((attempt + 1))
    done
    
    echo "Error: Emulator not ready on port $port after $max_attempts attempts"
    return 1
}

# Function to start emulator
start_emulator() {
    echo "Starting Spanner emulator..."
    cleanup_existing
    
    docker run -d --name spanner-emulator -p 9010:9010 -p 9020:9020 \
        gcr.io/cloud-spanner-emulator/emulator:latest
    
    # Check if container is running
    if ! docker ps --format 'table {{.Names}}' | grep -q spanner-emulator; then
        echo "Error: Emulator container is not running"
        docker logs spanner-emulator
        exit 1
    fi
    
    echo "Emulator container is running ✓"
    
    # Wait for emulator to be ready on both ports
    if ! waitForEmulator 9010; then
        echo "Error: Emulator failed to start on port 9010"
        docker logs spanner-emulator
        exit 1
    fi
    
    if ! waitForEmulator 9020; then
        echo "Error: Emulator failed to start on port 9020"
        docker logs spanner-emulator
        exit 1
    fi
}

# Main setup function
main() {
    check_dependencies
    
    # Start emulator if not running
    if ! docker ps --format 'table {{.Names}}' | grep -q spanner-emulator; then
        start_emulator
    else
        echo "Emulator already running ✓"
    fi
    
    # Set environment variables for emulator
    export SPANNER_EMULATOR_HOST=localhost:9010
    
    # Ensure emulator is fully ready before proceeding
    if ! waitForEmulator 9010; then
        echo "Error: Emulator is not ready"
        exit 1
    fi
    
    # Create instance and database using built-in schema initialization
    echo "Setting up instance and database..."
    cd "$(dirname "$0")/.."
    
    # Initialize database using built-in --init-schema functionality
    echo "Initializing database with built-in schema support..."
    go run . --project "$PROJECT_ID" --instance "$INSTANCE_ID" --database "$DATABASE_ID" --init-schema test/schema.sql --verbose || {
        echo "Failed to initialize database with --init-schema"
        exit 1
    }
    
    echo ""
    echo "🎉 Emulator setup complete!"
    echo "Project: $PROJECT_ID"
    echo "Instance: $INSTANCE_ID" 
    echo "Database: $DATABASE_ID"
    echo "Host: localhost:9010"
    echo ""
    echo "To run tests: make test-integration"
}

# Run main function
main "$@"