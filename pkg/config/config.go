package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SourceDatabase string            `yaml:"defaultSource"`
	TargetDatabase string            `yaml:"defaultTarget"`
	Stages         map[string]string `yaml:"stages"`
	Tables         []ConfigTable     `yaml:"tables"`
}

type ConfigTable struct {
	Name        string        `yaml:"name"`
	Columns     []string      `yaml:"columns"`
	PrimaryKeys []string      `yaml:"primaryKeys"`
	OnlyWhere   []ConfigWhere `yaml:"onlyWhere"`
	NoDelete    bool          `yaml:"noDelete"`
}

type ConfigWhere struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"` // string, bool, int, float, date
	Value string `yaml:"value"`
}

func ParseConfigFromFile(path string) (*Config, error) {
	// Read file from path

	f, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	// Parse file
	config := &Config{}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		return nil, err
	}

	// check for each table if it has a primary key
	for _, table := range config.Tables {
		if len(table.PrimaryKeys) == 0 {
			return nil, errors.New("table " + table.Name + " has no primary keys")
		}
	}

	return config, nil
}
