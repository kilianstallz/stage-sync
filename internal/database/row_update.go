package database

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"stage-sync-cli/internal/database/builder"
	"stage-sync-cli/models"
)

func UpdateRows(ctx context.Context, tx *sql.Tx, tableName string, changedColumns []string, oldRows []models.Row, updatedRows []models.Row, isDryRun bool) error {
	opt := zap.NewProductionConfig()
	opt.OutputPaths = []string{"update.log"}
	logger, err := opt.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	for _, oldRow := range oldRows {

		query := builder.BuildUpdateQuery(tableName, changedColumns, oldRow, updatedRows[0])
		logger.Sugar().Info(query)

		if !isDryRun {
			_, err := tx.ExecContext(ctx, query)
			if err != nil {
				panic(err)
			}

		}
	}

	return nil
}
