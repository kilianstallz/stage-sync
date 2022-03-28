package database

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"stage-sync-cli/internal/database/builder"
	"stage-sync-cli/models"
)

func DeleteRows(ctx context.Context, tx *sql.Tx, tableName string, rows []models.Row, isDryRun bool) error {
	for _, row := range rows {

		query := builder.BuildDeleteQuery(tableName, row)
		zap.S().Debug(query)

		if !isDryRun {
			_, err := tx.ExecContext(ctx, query)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}
