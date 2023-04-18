package postgres_test

import (
	"github.com/kilianstallz/stage-sync/internal/database/postgres"
	"github.com/kilianstallz/stage-sync/pkg/models"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBuildDeleteQuery(t *testing.T) {
	RegisterTestingT(t)
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

	query := postgres.BuildDeleteQuery(tableName, row)
	Expect(query).To(Equal("DELETE FROM \"users\" WHERE ((\"Id\" = 1) AND (\"Name\" = 'John') AND (\"Age\" = 1.6346346));"))

}
