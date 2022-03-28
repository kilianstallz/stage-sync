package builder_test

import (
	"fmt"
	"stage-sync-cli/config"
	"stage-sync-cli/internal/database/builder"
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

	query := builder.BuildSelectQuery(tableConfigs[0])

	if query != fmt.Sprintf("SELECT \"Id\", \"Name\", \"Age\" FROM %q WHERE \"Id\" = 1 AND \"Id\" = 2", tableConfigs[0].Name) {
		t.Errorf("SELECT \"Id\", \"Name\", \"Age\" FROM %q WHERE \"Id\" = 1 AND \"Id\" = 2, got '%s'", tableConfigs[0].Name, query)
	}

	query2 := builder.BuildSelectQuery(tableConfigs[1])

	if query2 != fmt.Sprintf("SELECT \"Name\", \"Age\" FROM %q", tableConfigs[1].Name) {
		t.Errorf("expected 'SELECT \"Name\", \"Age\" FROM %q', got '%s'", tableConfigs[1].Name, query2)
	}

}
