package atlas_test

import (
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"testing"
)

func TestInspect(t *testing.T) {

	ctx := context.Background()

	db, err := sql.Open("pgx", "postgres://localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	driver, err := postgres.Open(db)
	if err != nil {
		panic(err)
	}

	schem, err := driver.InspectSchema(ctx, "public", &schema.InspectOptions{})
	if err != nil {
		return
	}

	for _, table := range schem.Tables {

		primary := table.PrimaryKey
		if primary != nil {
			println("Primary key: " + primary.Name)
			for _, column := range primary.Parts {
				println("- " + column.C.Name)
			}
		}

		println(table.Name)
		for _, column := range table.Columns {
			str := "- " + column.Name + " "

			// add type
			str += column.Type.Raw

			if column.Type.Null {
				str += " NULL"
			} else {
				str += " NOT NULL"
			}

			println(str)
		}
	}

}
