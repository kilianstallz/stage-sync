package database

import (
	"context"
	"database/sql"
	"stage-sync/internal/database/builder"
	"stage-sync/internal/sql_log"
	"stage-sync/models"
)

func UpdateRows(ctx context.Context, tx *sql.Tx, tableName string, changedColumns []string, oldRows []models.Row, updatedRows []models.Row, isDryRun bool) error {
	log, f := sql_log.CreateSqlLogger("update.sql")
	defer f.Close()

	for i, oldRow := range oldRows {

		query := builder.BuildUpdateQuery(tableName, changedColumns, oldRow, updatedRows[i])
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
