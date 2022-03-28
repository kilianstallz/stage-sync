package builder_test

import (
	"fmt"
	"stage-sync-cli/internal/database/builder"
	"stage-sync-cli/models"
	"testing"
)

func TestBuildDeleteQuery(t *testing.T) {
	tableName := "users"
	row := models.Row{
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
	}

	query := builder.BuildDeleteQuery(tableName, row)
	if query != fmt.Sprintf("DELETE FROM %q WHERE \"Id\" = 1 AND \"Name\" = 'John' AND \"Age\" = 1.634635", tableName) {
		t.Errorf("expected 'DELETE FROM %q WHERE \"Id\" = 1 AND \"Name\" = 'John' AND \"Age\" = 1.634635', got '%s'", tableName, query)
	}
}
