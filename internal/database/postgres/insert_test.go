package postgres_test

import (
	"fmt"
	"github.com/kilianstallz/stage-sync/internal/database/postgres"
	"github.com/kilianstallz/stage-sync/pkg/models"
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

	query := postgres.BuildInsertQuery(tableName, rows[0])
	if query != fmt.Sprintf("INSERT INTO %q (\"Age\", \"Id\", \"Name\") VALUES (1.6346346, 1, 'John');", tableName) {
		t.Errorf("expected 'INSERT INTO %q (\"Age\", \"Id\", \"Name\") VALUES (1.6346346, 1, 'John');', got '%s'", tableName, query)
	}
}

func TestBuildInsertQueryWithIntArray(t *testing.T) {
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
				Name:  "Verbs",
				Value: []int{1, 2, 3},
				Type:  "[]int",
			},
		}}

	query := postgres.BuildInsertQuery(tableName, rows[0])
	if query != fmt.Sprintf("INSERT INTO %q (\"Id\", \"Name\", \"Verbs\") VALUES (1, 'John', (1, 2, 3));", tableName) {
		t.Errorf("expected 'INSERT INTO %q (\"Id\", \"Name\", \"Verbs\") VALUES (1, 'John', (1, 2, 3));', got '%s'", tableName, query)
	}
}
