package builder_test

import (
	"stage-sync-cli/config"
	"stage-sync-cli/internal/database/builder"
	"testing"
)

func TestBuildConnectionString(t *testing.T) {
	testCases := []config.ConfigDB{
		{
			User:     "user",
			Password: "password",
			Host:     "host",
			Port:     5432,
			Database: "database",
		},
		{
			User:     "",
			Password: "",
			Host:     "localhost",
			Port:     5432,
			Database: "database",
		},
	}
	results := []string{
		"postgres://user:password@host:5432/database?sslmode=require",
		"postgres://localhost:5432/database?sslmode=require",
	}

	for i, testCase := range testCases {
		result := builder.BuildConnectionString(testCase)
		if result != results[i] {
			t.Errorf("expected '%s', got '%s'", results[i], result)
		}
	}

}



