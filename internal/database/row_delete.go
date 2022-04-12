package database

import (
	"context"
	"database/sql"
	"stage-sync/internal/database/builder"
	"stage-sync/internal/sql_log"
	"stage-sync/models"
)

func DeleteRows(ctx context.Context, tx *sql.Tx, tableName string, rows []models.Row, isDryRun bool) error {
	log, f := sql_log.CreateSqlLogger("delete.sql")
	defer f.Close()

	for _, row := range rows {

		query := builder.BuildDeleteQuery(tableName, row)
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
