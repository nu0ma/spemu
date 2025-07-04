package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseDMLContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
		wantErr  bool
	}{
		{
			name:     "single INSERT statement",
			content:  `INSERT INTO users (id, name) VALUES (1, 'John');`,
			expected: []string{"INSERT INTO users (id, name) VALUES (1, 'John')"},
			wantErr:  false,
		},
		{
			name: "multiple statements",
			content: `INSERT INTO users (id, name) VALUES (1, 'John');
			UPDATE users SET name = 'Jane' WHERE id = 1;
			DELETE FROM users WHERE id = 2;`,
			expected: []string{
				"INSERT INTO users (id, name) VALUES (1, 'John')",
				"UPDATE users SET name = 'Jane' WHERE id = 1",
				"DELETE FROM users WHERE id = 2",
			},
			wantErr: false,
		},
		{
			name: "statements with comments",
			content: `-- This is a comment
			INSERT INTO users (id, name) VALUES (1, 'John'); -- inline comment
			-- Another comment
			UPDATE users SET name = 'Jane' WHERE id = 1;`,
			expected: []string{
				"INSERT INTO users (id, name) VALUES (1, 'John')",
				"UPDATE users SET name = 'Jane' WHERE id = 1",
			},
			wantErr: false,
		},
		{
			name: "empty content",
			content: `-- Only comments
			-- Nothing else`,
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "invalid statement",
			content:  `SELECT * FROM users;`,
			expected: nil,
			wantErr:  true,
		},
		{
			name: "mixed valid and invalid",
			content: `INSERT INTO users (id) VALUES (1);
			SELECT * FROM users;`,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDMLContent(tt.content)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseDMLContent() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseDMLContent() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseDMLContent() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestParseDMLFile(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.sql")

	content := `INSERT INTO users (id, name) VALUES (1, 'John');
	UPDATE users SET name = 'Jane' WHERE id = 1;`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ParseDMLFile(testFile)
	if err != nil {
		t.Errorf("ParseDMLFile() unexpected error: %v", err)
	}

	expected := []string{
		"INSERT INTO users (id, name) VALUES (1, 'John')",
		"UPDATE users SET name = 'Jane' WHERE id = 1",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseDMLFile() = %v, expected %v", result, expected)
	}
}

func TestParseDMLFileNotExist(t *testing.T) {
	_, err := ParseDMLFile("nonexistent.sql")
	if err == nil {
		t.Error("ParseDMLFile() expected error for nonexistent file")
	}
}

func TestRemoveComments(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "no comments",
			content:  "INSERT INTO users VALUES (1);",
			expected: "INSERT INTO users VALUES (1);",
		},
		{
			name:     "line comment",
			content:  "-- This is a comment\nINSERT INTO users VALUES (1);",
			expected: "\nINSERT INTO users VALUES (1);",
		},
		{
			name:     "inline comment",
			content:  "INSERT INTO users VALUES (1); -- comment",
			expected: "INSERT INTO users VALUES (1); ",
		},
		{
			name:     "multiple comments",
			content:  "-- Comment 1\nINSERT INTO users VALUES (1); -- Comment 2\n-- Comment 3",
			expected: "\nINSERT INTO users VALUES (1); \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeComments(tt.content)
			if result != tt.expected {
				t.Errorf("removeComments() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestSplitStatements(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "single statement",
			content:  "INSERT INTO users VALUES (1)",
			expected: []string{"INSERT INTO users VALUES (1)"},
		},
		{
			name:     "multiple statements",
			content:  "INSERT INTO users VALUES (1); UPDATE users SET name = 'test';",
			expected: []string{"INSERT INTO users VALUES (1)", "UPDATE users SET name = 'test'"},
		},
		{
			name:     "statements with whitespace",
			content:  "  INSERT INTO users VALUES (1);  \n  UPDATE users SET name = 'test';  ",
			expected: []string{"INSERT INTO users VALUES (1)", "UPDATE users SET name = 'test'"},
		},
		{
			name:     "empty statements",
			content:  "INSERT INTO users VALUES (1);;;UPDATE users SET name = 'test';",
			expected: []string{"INSERT INTO users VALUES (1)", "UPDATE users SET name = 'test'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitStatements(tt.content)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("splitStatements() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsValidDMLStatement(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		expected  bool
	}{
		{"INSERT uppercase", "INSERT INTO users VALUES (1)", true},
		{"INSERT lowercase", "insert into users values (1)", true},
		{"INSERT mixed case", "Insert Into users Values (1)", true},
		{"UPDATE uppercase", "UPDATE users SET name = 'test'", true},
		{"UPDATE lowercase", "update users set name = 'test'", true},
		{"DELETE uppercase", "DELETE FROM users WHERE id = 1", true},
		{"DELETE lowercase", "delete from users where id = 1", true},
		{"SELECT statement", "SELECT * FROM users", false},
		{"CREATE statement", "CREATE TABLE users (id INT)", false},
		{"DROP statement", "DROP TABLE users", false},
		{"Empty statement", "", false},
		{"Whitespace only", "   ", false},
		{"Invalid prefix", "INVALID STATEMENT", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidDMLStatement(tt.statement)
			if result != tt.expected {
				t.Errorf("isValidDMLStatement(%q) = %v, expected %v", tt.statement, result, tt.expected)
			}
		})
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
