package config

import "fmt"

type Config struct {
	EmulatorHost string
	ProjectID    string
	InstanceID   string
	DatabaseID   string
}

func (c *Config) DatabasePath() string {
	return fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		c.ProjectID, c.InstanceID, c.DatabaseID)
}