package test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/nu0ma/spemu/pkg/config"
	"github.com/nu0ma/spemu/pkg/executor"
	"github.com/nu0ma/spemu/pkg/parser"
)

const (
	testProjectID  = "test-project"
	testInstanceID = "test-instance"
	testDatabaseID = "test-database"
	emulatorHost   = "localhost:9010"
)

func TestMain(m *testing.M) {
	// Setup: Check if emulator is available
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		os.Setenv("SPANNER_EMULATOR_HOST", emulatorHost)
	}

	// Setup emulator instance and database
	if err := setupEmulator(); err != nil {
		fmt.Printf("Failed to setup emulator: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func setupEmulator() error {
	// Check if setup script exists (look from project root)
	scriptPath := "../scripts/setup-emulator.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		// Try from current directory (if running from project root)
		scriptPath = "scripts/setup-emulator.sh"
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return fmt.Errorf("setup script not found: %s", scriptPath)
		}
	}

	// Run setup script
	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func setupTestDatabase(t *testing.T) (*spanner.Client, func()) {
	t.Helper()

	ctx := context.Background()
	cfg := &config.Config{
		ProjectID:    testProjectID,
		InstanceID:   testInstanceID,
		DatabaseID:   testDatabaseID,
		EmulatorHost: emulatorHost,
	}

	client, err := spanner.NewClient(ctx, cfg.DatabasePath())
	if err != nil {
		t.Skipf("Failed to create Spanner client (emulator may not be running): %v", err)
	}

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		_, err := client.Apply(ctx, []*spanner.Mutation{
			spanner.Delete("comments", spanner.AllKeys()),
			spanner.Delete("posts", spanner.AllKeys()),
			spanner.Delete("users", spanner.AllKeys()),
			spanner.Delete("test_table", spanner.AllKeys()),
		})
		if err != nil {
			t.Logf("Failed to cleanup test data: %v", err)
		}
		client.Close()
	}

	return client, cleanup
}

func TestIntegration_ExecuteStatements(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client, cleanup := setupTestDatabase(t)
	defer cleanup()

	cfg := &config.Config{
		ProjectID:    testProjectID,
		InstanceID:   testInstanceID,
		DatabaseID:   testDatabaseID,
		EmulatorHost: emulatorHost,
	}

	exec, err := executor.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}
	defer exec.Close()

	// Test data insertion
	statements := []string{
		"INSERT INTO users (id, name, email, created_at) VALUES (1, 'John Doe', 'john@example.com', '2024-01-01T00:00:00Z')",
		"INSERT INTO users (id, name, email, created_at) VALUES (2, 'Jane Smith', 'jane@example.com', '2024-01-02T00:00:00Z')",
		"INSERT INTO posts (id, user_id, title, content, created_at) VALUES (1, 1, 'Test Post', 'This is a test post', '2024-01-01T01:00:00Z')",
	}

	err = exec.ExecuteStatements(statements, true)
	if err != nil {
		t.Fatalf("Failed to execute statements: %v", err)
	}

	// Verify data was inserted
	ctx := context.Background()

	// Check users table
	iter := client.Single().Query(ctx, spanner.Statement{SQL: "SELECT COUNT(*) as count FROM users"})
	defer iter.Stop()
	row, err := iter.Next()
	if err != nil {
		t.Fatalf("Failed to query users count: %v", err)
	}
	var userCount int64
	if err := row.Columns(&userCount); err != nil {
		t.Fatalf("Failed to scan user count: %v", err)
	}
	if userCount != 2 {
		t.Errorf("Expected 2 users, got %d", userCount)
	}

	// Check posts table
	iter = client.Single().Query(ctx, spanner.Statement{SQL: "SELECT COUNT(*) as count FROM posts"})
	defer iter.Stop()
	row, err = iter.Next()
	if err != nil {
		t.Fatalf("Failed to query posts count: %v", err)
	}
	var postCount int64
	if err := row.Columns(&postCount); err != nil {
		t.Fatalf("Failed to scan post count: %v", err)
	}
	if postCount != 1 {
		t.Errorf("Expected 1 post, got %d", postCount)
	}
}

func TestIntegration_ParseAndExecute(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	_, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Test the complete flow: parse SQL file and execute
	sqlContent := `-- Test data
INSERT INTO users (id, name, email, created_at) VALUES (10, 'Test User', 'test@example.com', '2024-01-01T00:00:00Z');
-- Insert a post
INSERT INTO posts (id, user_id, title, content, created_at) VALUES (10, 10, 'Integration Test', 'Testing the complete flow', '2024-01-01T01:00:00Z');`

	statements, err := parser.ParseDMLContent(sqlContent)
	if err != nil {
		t.Fatalf("Failed to parse SQL content: %v", err)
	}

	if len(statements) != 2 {
		t.Errorf("Expected 2 statements, got %d", len(statements))
	}

	cfg := &config.Config{
		ProjectID:    testProjectID,
		InstanceID:   testInstanceID,
		DatabaseID:   testDatabaseID,
		EmulatorHost: emulatorHost,
	}

	exec, err := executor.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}
	defer exec.Close()

	err = exec.ExecuteStatements(statements, false)
	if err != nil {
		t.Fatalf("Failed to execute parsed statements: %v", err)
	}
}

func TestIntegration_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	_, cleanup := setupTestDatabase(t)
	defer cleanup()

	cfg := &config.Config{
		ProjectID:    testProjectID,
		InstanceID:   testInstanceID,
		DatabaseID:   testDatabaseID,
		EmulatorHost: emulatorHost,
	}

	exec, err := executor.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}
	defer exec.Close()

	// Test with invalid SQL (should fail)
	invalidStatements := []string{
		"INSERT INTO nonexistent_table (id) VALUES (1)",
	}

	err = exec.ExecuteStatements(invalidStatements, false)
	if err == nil {
		t.Error("Expected error when inserting into nonexistent table")
	}

	// Test with constraint violation (duplicate primary key)
	statements := []string{
		"INSERT INTO users (id, name, email, created_at) VALUES (100, 'User 1', 'user1@example.com', '2024-01-01T00:00:00Z')",
		"INSERT INTO users (id, name, email, created_at) VALUES (100, 'User 2', 'user2@example.com', '2024-01-01T00:00:00Z')", // Same ID
	}

	err = exec.ExecuteStatements(statements, false)
	if err == nil {
		t.Error("Expected error when inserting duplicate primary key")
	}
}

func TestIntegration_TransactionRollback(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client, cleanup := setupTestDatabase(t)
	defer cleanup()

	cfg := &config.Config{
		ProjectID:    testProjectID,
		InstanceID:   testInstanceID,
		DatabaseID:   testDatabaseID,
		EmulatorHost: emulatorHost,
	}

	exec, err := executor.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}
	defer exec.Close()

	// First, insert a valid user
	validStatements := []string{
		"INSERT INTO users (id, name, email, created_at) VALUES (200, 'Valid User', 'valid@example.com', '2024-01-01T00:00:00Z')",
	}
	err = exec.ExecuteStatements(validStatements, false)
	if err != nil {
		t.Fatalf("Failed to insert valid user: %v", err)
	}

	// Now try a transaction that should fail (and rollback)
	mixedStatements := []string{
		"INSERT INTO users (id, name, email, created_at) VALUES (201, 'User 201', 'user201@example.com', '2024-01-01T00:00:00Z')",
		"INSERT INTO nonexistent_table (id) VALUES (1)", // This will fail
		"INSERT INTO users (id, name, email, created_at) VALUES (202, 'User 202', 'user202@example.com', '2024-01-01T00:00:00Z')",
	}

	err = exec.ExecuteStatements(mixedStatements, false)
	if err == nil {
		t.Error("Expected error when executing statements with invalid table")
	}

	// Verify that no users were inserted (transaction should have rolled back)
	ctx := context.Background()
	iter := client.Single().Query(ctx, spanner.Statement{SQL: "SELECT COUNT(*) as count FROM users WHERE id IN (201, 202)"})
	defer iter.Stop()
	row, err := iter.Next()
	if err != nil {
		t.Fatalf("Failed to query user count: %v", err)
	}
	var count int64
	if err := row.Columns(&count); err != nil {
		t.Fatalf("Failed to scan count: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 users from failed transaction, got %d", count)
	}
}

func TestIntegration_PerformanceWithLargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	_, cleanup := setupTestDatabase(t)
	defer cleanup()

	cfg := &config.Config{
		ProjectID:    testProjectID,
		InstanceID:   testInstanceID,
		DatabaseID:   testDatabaseID,
		EmulatorHost: emulatorHost,
	}

	exec, err := executor.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}
	defer exec.Close()

	// Generate multiple insert statements
	var statements []string
	for i := 1; i <= 100; i++ {
		statements = append(statements,
			fmt.Sprintf("INSERT INTO test_table (id, name, value, created_at) VALUES (%d, 'Test User %d', 'Value %d', '2024-01-01T00:00:00Z')", i, i, i))
	}

	start := time.Now()
	err = exec.ExecuteStatements(statements, false)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to execute large dataset: %v", err)
	}

	t.Logf("Inserted 100 records in %v", duration)

	// Performance expectation: should complete within reasonable time
	if duration > 30*time.Second {
		t.Errorf("Large dataset insertion took too long: %v", duration)
	}
}
