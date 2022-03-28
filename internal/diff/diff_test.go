package diff_test

import (
	"stage-sync-cli/internal/diff"
	"stage-sync-cli/models"
	"testing"
)

func TestFindDiffResult(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"basicRowUpdate": testBasicRowUpdate,
		"basicRowDelete": testBasicRowDelete,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

func testBasicRowUpdate(t *testing.T) {
	bTable := models.Table{
		Name: "Base",
		PrimaryKeys: []string{
			"Id",
		},
		Rows: []models.Row{
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "1",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "John",
				},
			},
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "2",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "Jane",
				},
			},
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "3",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "Jack",
				},
			},
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "4",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "Jack",
				},
			},
		},
	}

	dTable := models.Table{
		Name: "Diff",
		PrimaryKeys: []string{
			"Id",
		},
		Rows: []models.Row{
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "1",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "John",
				},
			},
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "2",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "Janny",
				},
			},
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "3",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "Jack",
				},
			},
		},
	}

	diffRes := diff.FindDiffResult(bTable, dTable)

	if len(diffRes.UpdatedRows.ChangedColumns) != 1 {
		t.Errorf("Expected 1 updated row, got %d", len(diffRes.UpdatedRows.ChangedColumns))
	}

	if len(diffRes.AddedRows) != 1 {
		t.Errorf("Expected 1 added row, got %d", len(diffRes.UpdatedRows.ChangedColumns))
	}
}

func testBasicRowDelete(t *testing.T) {
	bTable := models.Table{
		Name: "Base",
		PrimaryKeys: []string{
			"Id",
		},
		Rows: []models.Row{
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "1",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "John",
				},
			},
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "2",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "Jane",
				},
			},
		},
	}

	dTable := models.Table{
		Name: "Diff",
		PrimaryKeys: []string{
			"Id",
		},
		Rows: []models.Row{
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "1",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "John",
				},
			},
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "2",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "Janny",
				},
			},
			[]models.Column{
				models.Column{
					Name:  "Id",
					Type:  "int",
					Value: "3",
				},
				models.Column{
					Name:  "Name",
					Type:  "string",
					Value: "Jack",
				},
			},
		},
	}

	diffRes := diff.FindDiffResult(bTable, dTable)

	if len(diffRes.DeletedRows) != 1 {
		t.Errorf("Expected 1 deleted row, got %d", len(diffRes.DeletedRows))
	}
}