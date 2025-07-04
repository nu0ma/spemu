package parser

import (
	"fmt"
	"os"
	"strings"
)

func ParseDMLFile(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return ParseDMLContent(string(content))
}

func ParseDMLContent(content string) ([]string, error) {
	content = removeComments(content)

	statements := splitStatements(content)

	var validStatements []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if !isValidDMLStatement(stmt) {
			limit := 50
			if len(stmt) < limit {
				limit = len(stmt)
			}
			return nil, fmt.Errorf("invalid DML statement: %s", stmt[:limit])
		}

		validStatements = append(validStatements, stmt)
	}

	return validStatements, nil
}

func removeComments(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		if idx := strings.Index(line, "--"); idx != -1 {
			line = line[:idx]
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func splitStatements(content string) []string {
	statements := strings.Split(content, ";")

	var result []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	return result
}

func isValidDMLStatement(stmt string) bool {
	stmt = strings.TrimSpace(strings.ToUpper(stmt))

	validPrefixes := []string{
		"INSERT",
		"UPDATE",
		"DELETE",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(stmt, prefix) {
			return true
		}
	}

	return false
}

