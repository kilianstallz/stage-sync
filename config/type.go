package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	SourceDatabase DbConnection `yaml:"sourceDatabase"`
	TargetDatabase DbConnection `yaml:"targetDatabase"`
	Tables []ConfigTable `yaml:"tables"`
}

type DbConnection struct {
	Credentials ConfigDB `yaml:"credentials"`
	Schema      string   `yaml:"schema"`
}

type ConfigTable struct {
	Name        string        `yaml:"name"`
	Columns     []string      `yaml:"columns"`
	PrimaryKeys []string      `yaml:"primaryKeys"`
	OnlyWhere   []ConfigWhere `yaml:"onlyWhere"`
	NoDelete 	bool          `yaml:"noDelete"`
}

type ConfigWhere struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"` // string, bool, int, float, date
	Value string `yaml:"value"`
}

type ConfigDB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func ParseConfigFromFile(path string) (*Config, error) {
	// Read file from path
	f, err := ioutil.ReadFile(path)
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
