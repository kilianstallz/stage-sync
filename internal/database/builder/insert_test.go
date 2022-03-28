package builder_test

import (
	"fmt"
	"stage-sync-cli/internal/database/builder"
	"stage-sync-cli/models"
	"testing"
)

func TestBuildInsertQuery(t *testing.T) {
	tableName := "users"
	rows := []models.Row{
		{
			models.Column{
				Name:  "Id",
				Value: 1,
				Type:  "int",
			},
			models.Column{
				Name:  "Name",
				Value: "John",
				Type:  "string",
			},
			models.Column{
				Name:  "Age",
				Value: 1.6346346,
				Type:  "float64",
			},
		},
	}

	query := builder.BuildInsertQuery(tableName, rows[0])
	if query != fmt.Sprintf("INSERT INTO %q (\"Id\", \"Name\", \"Age\") VALUES (1, 'John', 1.6346346)", tableName) {
		t.Errorf("expected 'INSERT INTO %q (\"Id\", \"Name\", \"Age\") VALUES (1, 'John', 1.6346346)', got '%s'", tableName, query)
	}
}
