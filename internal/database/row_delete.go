package database

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"stage-sync-cli/internal/database/builder"
	"stage-sync-cli/models"
)

func DeleteRows(ctx context.Context, tx *sql.Tx, tableName string, rows []models.Row, isDryRun bool) error {
	opt := zap.NewProductionConfig()
	opt.OutputPaths = []string{"delete.log"}
	logger, err := opt.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	for _, row := range rows {

		query := builder.BuildDeleteQuery(tableName, row)
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
