package postgres_test

import (
	"fmt"
	"github.com/kilianstallz/stage-sync/internal/database/postgres"
	"github.com/kilianstallz/stage-sync/models"
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
		models.Column{
			Name:  "Text",
			Value: "<html> this is some \"text\" </html>",
			Type:  "string",
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
		models.Column{
			Name:  "Text",
			Value: "<html> this is some \"new text\" </html>",
			Type:  "string",
		},
	}
	changedColumns := []string{
		"Age",
		"Text",
	}

	query := postgres.BuildUpdateQuery(tableName, changedColumns, oldRow, newRow)
	t.Log(query)
	if query != fmt.Sprintf("UPDATE %q SET \"Text\"='<html> this is some \"new text\" </html>' WHERE ((\"Id\" = 1) AND (\"Name\" = 'John') AND (\"Age\" = 1.6346346) AND (\"Text\" = '<html> this is some \"text\" </html>'));", tableName) {
		t.Errorf("UPDATE %q SET \"Text\"='<html> this is some \"new text\" </html>' WHERE ((\"Id\" = 1) AND (\"Name\" = 'John') AND (\"Age\" = 1.6346346) AND (\"Text\" = '<html> this is some \"text\" </html>'));' \n got \n '%s'", tableName, query)
	}

}
