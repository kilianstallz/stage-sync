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

func DeleteRows(ctx context.Context, tx *sql.Tx, tableName string, rows []models.Row, isDryRun bool) error {
	f, err := os.OpenFile("delete.sql", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		zap.S().Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetFlags(0)
	log.SetOutput(f)

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
