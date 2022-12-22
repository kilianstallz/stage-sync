package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	SourceDatabase string        `yaml:"sourceDatabase"`
	TargetDatabase string        `yaml:"targetDatabase"`
	Tables         []ConfigTable `yaml:"tables"`
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
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Parse file
	config := &Config{}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
