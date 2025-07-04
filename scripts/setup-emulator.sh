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
    
    if ! command -v gcloud &> /dev/null; then
        echo "Error: gcloud CLI is required but not installed"
        echo "Please install Google Cloud SDK: https://cloud.google.com/sdk/docs/install"
        exit 1
    fi
    
    if ! command -v curl &> /dev/null; then
        echo "Error: curl is required but not installed"
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

# Function to start emulator
start_emulator() {
    echo "Starting Spanner emulator..."
    cleanup_existing
    
    docker run -d --name spanner-emulator -p 9010:9010 -p 9020:9020 \
        gcr.io/cloud-spanner-emulator/emulator:latest
    
    # Wait for emulator to be ready
    echo "Waiting for emulator to be ready..."
    sleep 5  # Give it time to start
    
    # Check if container is running
    if ! docker ps --format 'table {{.Names}}' | grep -q spanner-emulator; then
        echo "Error: Emulator container is not running"
        docker logs spanner-emulator
        exit 1
    fi
    
    echo "Emulator container is running ✓"
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
    
    # Wait a bit more for the emulator to be fully ready
    echo "Waiting for emulator to be fully ready..."
    sleep 3
    
    # Create instance and database using Go script
    echo "Setting up instance and database..."
    cd "$(dirname "$0")/.."
    go run scripts/setup.go
    
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