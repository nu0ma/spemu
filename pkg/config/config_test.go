package config

import (
	"testing"
)

func TestConfig_DatabasePath(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "basic configuration",
			config: Config{
				ProjectID:  "test-project",
				InstanceID: "test-instance",
				DatabaseID: "test-database",
			},
			expected: "projects/test-project/instances/test-instance/databases/test-database",
		},
		{
			name: "configuration with special characters",
			config: Config{
				ProjectID:  "my-test-project-123",
				InstanceID: "my-instance-456",
				DatabaseID: "my-database-789",
			},
			expected: "projects/my-test-project-123/instances/my-instance-456/databases/my-database-789",
		},
		{
			name: "configuration with emulator host (should not affect path)",
			config: Config{
				EmulatorHost: "localhost:9010",
				ProjectID:    "emulator-project",
				InstanceID:   "emulator-instance",
				DatabaseID:   "emulator-database",
			},
			expected: "projects/emulator-project/instances/emulator-instance/databases/emulator-database",
		},
		{
			name: "empty values",
			config: Config{
				ProjectID:  "",
				InstanceID: "",
				DatabaseID: "",
			},
			expected: "projects//instances//databases/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.DatabasePath()
			if result != tt.expected {
				t.Errorf("Config.DatabasePath() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestConfig_DatabasePath_Consistency(t *testing.T) {
	config := Config{
		ProjectID:  "test-project",
		InstanceID: "test-instance",
		DatabaseID: "test-database",
	}

	// Call multiple times to ensure consistency
	path1 := config.DatabasePath()
	path2 := config.DatabasePath()
	path3 := config.DatabasePath()

	if path1 != path2 || path2 != path3 {
		t.Errorf("DatabasePath() should return consistent results: %q, %q, %q", path1, path2, path3)
	}
}

func TestConfig_Fields(t *testing.T) {
	config := Config{
		EmulatorHost: "localhost:9010",
		ProjectID:    "test-project",
		InstanceID:   "test-instance",
		DatabaseID:   "test-database",
	}

	// Test that all fields are accessible
	if config.EmulatorHost != "localhost:9010" {
		t.Errorf("EmulatorHost = %q, expected %q", config.EmulatorHost, "localhost:9010")
	}
	if config.ProjectID != "test-project" {
		t.Errorf("ProjectID = %q, expected %q", config.ProjectID, "test-project")
	}
	if config.InstanceID != "test-instance" {
		t.Errorf("InstanceID = %q, expected %q", config.InstanceID, "test-instance")
	}
	if config.DatabaseID != "test-database" {
		t.Errorf("DatabaseID = %q, expected %q", config.DatabaseID, "test-database")
	}
}