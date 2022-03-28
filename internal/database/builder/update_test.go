package builder_test

import (
	"fmt"
	"stage-sync-cli/internal/database/builder"
	"stage-sync-cli/models"
	"testing"
)

func TestBuildUpdateQuery(t *testing.T) {
	tableName := "users"
	oldRow := models.Row{
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
	newRow := models.Row{
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
			Value: 2.5,
			Type:  "float64",
		},
	}
	changedColumns := []string{
		"Age",
	}

	query := builder.BuildUpdateQuery(tableName, changedColumns, oldRow, newRow)
	t.Log(query)
	if query != fmt.Sprintf("UPDATE %q SET \"Age\" = 2.500000 WHERE \"Id\" = 1 AND \"Name\" = 'John' AND \"Age\" = 1.634635", tableName) {
		t.Errorf("expected \n 'UPDATE %q SET \"Age\" = 2.5 WHERE \"Id\" = 1 AND \"Name\" = 'John' AND \"Age\" = 1.634635' \n got \n '%s'", tableName, query)
	}

}
