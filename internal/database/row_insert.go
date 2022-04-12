package database

import (
	"context"
	"database/sql"
	"stage-sync/internal/database/builder"
	"stage-sync/internal/sql_log"
	"stage-sync/models"
)

func InsertRows(ctx context.Context, tx *sql.Tx, tableName string, rows []models.Row, isDryRun bool) error {
	log, f := sql_log.CreateSqlLogger("insert.sql")
	defer f.Close()

	for _, row := range rows {

		query := builder.BuildInsertQuery(tableName, row)
		log.Println(query)

		if !isDryRun {
			_, err := tx.ExecContext(ctx, query)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}
