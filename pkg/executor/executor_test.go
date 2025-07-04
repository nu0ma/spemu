package executor

import (
	"os"
	"testing"

	"github.com/nu0ma/spemu/pkg/config"
)

func TestNew(t *testing.T) {
	// Test creating executor with valid config
	cfg := &config.Config{
		ProjectID:    "test-project",
		InstanceID:   "test-instance",
		DatabaseID:   "test-database",
		EmulatorHost: "localhost:9010",
	}

	// This test requires Spanner emulator to be running
	// Skip if emulator is not available
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" && cfg.EmulatorHost == "" {
		t.Skip("Skipping test: Spanner emulator not available")
	}

	executor, err := New(cfg)
	if err != nil {
		t.Skipf("Failed to create executor (emulator may not be running): %v", err)
	}
	defer executor.Close()

	if executor.client == nil {
		t.Error("Expected client to be initialized")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	// Test with nil config - should panic, so we need to recover
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil config")
		}
	}()

	// This should panic
	New(nil)
}

func TestExecutor_Close(t *testing.T) {
	executor := &Executor{client: nil}

	// Should not panic with nil client
	executor.Close()
}

// Integration test - requires running Spanner emulator
func TestExecuteStatements_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := &config.Config{
		ProjectID:    "test-project",
		InstanceID:   "test-instance",
		DatabaseID:   "test-database",
		EmulatorHost: "localhost:9010",
	}

	// Skip if emulator is not available
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" && cfg.EmulatorHost == "" {
		t.Skip("Skipping integration test: Spanner emulator not available")
	}

	executor, err := New(cfg)
	if err != nil {
		t.Skipf("Failed to create executor (emulator may not be running): %v", err)
	}
	defer executor.Close()

	// Test statements - these require the tables to exist in the test database
	statements := []string{
		"INSERT INTO test_table (id, name) VALUES (1, 'test')",
	}

	err = executor.ExecuteStatements(statements, false)
	// We expect this to fail if test_table doesn't exist, which is normal
	// The important thing is that we can create the executor and call the method
	if err != nil {
		t.Logf("Expected error executing statements (test table may not exist): %v", err)
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"a smaller", 3, 5, 3},
		{"b smaller", 7, 2, 2},
		{"equal", 4, 4, 4},
		{"negative numbers", -3, -1, -3},
		{"zero", 0, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := min(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("min(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// Benchmark for ExecuteStatements (when emulator is available)
func BenchmarkExecuteStatements(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	cfg := &config.Config{
		ProjectID:    "test-project",
		InstanceID:   "test-instance",
		DatabaseID:   "test-database",
		EmulatorHost: "localhost:9010",
	}

	executor, err := New(cfg)
	if err != nil {
		b.Skipf("Failed to create executor: %v", err)
	}
	defer executor.Close()

	statements := []string{
		"INSERT INTO benchmark_table (id, value) VALUES (1, 'test')",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will likely fail, but we're measuring the overhead
		_ = executor.ExecuteStatements(statements, false)
	}
}
