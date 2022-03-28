package propagate

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"stage-sync-cli/config"
	"stage-sync-cli/internal/database"
	"stage-sync-cli/internal/database/builder"
	"stage-sync-cli/internal/diff"
	"stage-sync-cli/internal/table"
)

func Run(configPath string, isDryRun bool) {
	zap.L().Info("Starting propagation")
	conf, _ := config.ParseConfigFromFile(configPath)

	// Get the source database connection string.
	sourceDbConnectionString := builder.BuildConnectionString(conf.SourceDatabase.Credentials)

	// Get the source database connection.
	sourceDB, err := database.NewPostgresClient(sourceDbConnectionString)
	if err != nil {
		zap.S().Fatal("Failed to connect to source database", zap.Error(err))
	}

	tables := database.QueryTables(conf, sourceDB)

	sourceDB.Close()

	// Get the target database connection string.
	targetDbConnectionString := builder.BuildConnectionString(conf.TargetDatabase.Credentials)

	// Get the target database connection.
	targetDB, err := database.NewPostgresClient(targetDbConnectionString)
	if err != nil {
		zap.S().Fatal("Failed to connect to target database", zap.Error(err))
	}

	targetTables := database.QueryTables(conf, targetDB)

	defer func(targetDB *sql.DB) {
		err := targetDB.Close()
		if err != nil {
			zap.S().Error("Failed to close target database connection", zap.Error(err))
		}
	}(targetDB)

	ctx := context.Background()

	tx, err := targetDB.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	// diff the data for each table based on the list of primary keys
	for _, ftable := range tables {

		// find the target table with the same name
		targetTable := table.FindTable(targetTables, ftable.Name)

		diffResult := diff.FindDiffResult(ftable, targetTable)

		zap.L().Info(fmt.Sprintf("%d added rows", len(diffResult.AddedRows)))
		zap.L().Info(fmt.Sprintf("%d deleted rows", len(diffResult.DeletedRows)))
		zap.L().Info(fmt.Sprintf("%d updated rows", len(diffResult.UpdatedRows.ChangedColumns)))

		if len(diffResult.AddedRows) == 0 && len(diffResult.DeletedRows) == 0 && len(diffResult.UpdatedRows.ChangedColumns) == 0 {
			zap.L().Info(fmt.Sprintf("No changes on Table: %s", ftable.Name))
			continue
		}

		if len(diffResult.AddedRows) > 0 {
			err := database.InsertRows(ctx, tx, targetTable.Name, diffResult.AddedRows, isDryRun)
			if err != nil {
				zap.L().Error("Error inserting rows ", zap.Error(err))
				return
			}
			zap.L().Info(fmt.Sprintf("Insert done on Table: %s", ftable.Name))

		}

		if len(diffResult.DeletedRows) > 0 {
			// delete the deleted rows
			err := database.DeleteRows(ctx, tx, targetTable.Name, diffResult.DeletedRows, isDryRun)
			if err != nil {
				zap.L().Error("Error deleting rows ", zap.Error(err))
				return
			}

			zap.L().Info(fmt.Sprintf("Delete done on Table: %s", ftable.Name))
		}

		if len(diffResult.UpdatedRows.ChangedColumns) > 0 {
			// update the updated rows
			err = database.UpdateRows(ctx, tx, targetTable.Name, diffResult.UpdatedRows.ChangedColumns, diffResult.UpdatedRows.Before, diffResult.UpdatedRows.After, isDryRun)
			if err != nil {
				zap.S().Error("Error updating rows ", zap.Error(err))
				return
			}
			zap.L().Info(fmt.Sprintf("Update done on Table: %s", ftable.Name))
		}
	}

	if !isDryRun {
		// commit the transaction
		err = tx.Commit()
		if err != nil {
			zap.S().Error("Failed to commit transaction", zap.Error(err))
			return
		}
	} else {
		zap.L().Info("Dry run, skipping transaction commit")
		err = tx.Rollback()
		if err != nil {
			zap.S().Error("Failed to rollback transaction", zap.Error(err))
			return
		}
	}
}
