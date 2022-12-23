package postgres_test

import (
	"fmt"
	"github.com/kilianstallz/stage-sync/internal/database/postgres"
	"github.com/kilianstallz/stage-sync/pkg/config"
	"testing"
)

func TestBuildSelectQuery(t *testing.T) {
	tableConfigs := []config.ConfigTable{
		{
			Name: "table1",
			PrimaryKeys: []string{
				"Id",
			},
			Columns: []string{
				"Id",
				"Name",
				"Age",
			},
			OnlyWhere: []config.ConfigWhere{
				{
					Name:  "Id",
					Value: "1",
					Type:  "int",
				},
				{
					Name:  "Id",
					Value: "2",
					Type:  "int",
				},
			},
		},
	}

	query := postgres.BuildSelectQuery(tableConfigs[0])
	t.Log(query)

	if query != fmt.Sprintf("SELECT \"Id\", \"Name\", \"Age\" FROM \"table1\" WHERE (\"Id\" IN ('1', '2'));") {
		t.Errorf("'SELECT \"Id\", \"Name\", \"Age\" FROM \"table1\" WHERE (\"Id\" IN ('1', '2'));', got '%s'", query)
	}
}

func Test_SingleSelect(t *testing.T) {
	tableConfigs := []config.ConfigTable{
		{
			Name: "table2",
			PrimaryKeys: []string{
				"Name",
			},
			Columns: []string{
				"Name",
				"Age",
			},
		},
	}
	query2 := postgres.BuildSelectQuery(tableConfigs[0])
	t.Log(query2)

	if query2 != fmt.Sprintf("SELECT \"Name\", \"Age\" FROM \"table2\";") {
		t.Errorf("expected 'SELECT \"Name\", \"Age\" FROM \"table2\";', got '%s'", query2)
	}
}
