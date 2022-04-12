package database

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"log"
	"os"
	"stage-sync-cli/internal/database/builder"
	"stage-sync-cli/models"
)

func UpdateRows(ctx context.Context, tx *sql.Tx, tableName string, changedColumns []string, oldRows []models.Row, updatedRows []models.Row, isDryRun bool) error {
	f, err := os.OpenFile("update.sql", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		zap.S().Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetFlags(0)
	log.SetOutput(f)

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
