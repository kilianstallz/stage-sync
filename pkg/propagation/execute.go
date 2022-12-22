package propagation

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"stage-sync/internal/config"
	"stage-sync/internal/database/builder"
	"stage-sync/internal/diff"
	"stage-sync/internal/table"
)

func Execute(configPath string, isDryRun bool) {
	zap.L().Info("Starting propagation")
	conf, _ := config.ParseConfigFromFile(configPath)

	// Get the source database connection.
	sourceDB := builder.NewPostgresClient(conf.SourceDatabase)

	tables := sourceDB.QueryTables(conf)
	sourceDB.Close()

	// Get the target database connection.
	targetDB := builder.NewPostgresClient(conf.TargetDatabase)

	targetTables := targetDB.QueryTables(conf)

	defer func(targetDB builder.QueryBuilder) {
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

	// recover from panic and rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			zap.S().Fatal("Recovered from panic with rollback", zap.Any("panic", r))
		}
	}()

	// diff the data for each table based on the list of primary keys
	for _, ftable := range tables {

		// find the target table with the same name
		targetTable := table.FindTable(targetTables, ftable.Name)

		diffResult := diff.FindDiffResult(ftable, targetTable)

		zap.L().Info(fmt.Sprintf("%d added rows", len(diffResult.AddedRows)))
		zap.L().Info(fmt.Sprintf("%d updated rows", len(diffResult.UpdatedRows.ChangedColumns)))
		if ftable.NoDelete == false {
			zap.L().Info(fmt.Sprintf("%d deleted rows", len(diffResult.DeletedRows)))
		}

		if len(diffResult.AddedRows) == 0 && len(diffResult.DeletedRows) == 0 && len(diffResult.UpdatedRows.ChangedColumns) == 0 {
			zap.L().Info(fmt.Sprintf("No changes on Table: %s", ftable.Name))
			continue
		}

		if len(diffResult.AddedRows) > 0 {
			err := targetDB.InsertRows(ctx, tx, targetTable.Name, diffResult.AddedRows, isDryRun)
			if err != nil {
				zap.L().Error("Error inserting rows ", zap.Error(err))
				return
			}
			zap.L().Info(fmt.Sprintf("Insert done on Table: %s", ftable.Name))

		}

		if len(diffResult.DeletedRows) > 0 {
			if ftable.NoDelete {
				zap.L().Info(fmt.Sprintf("No delete on Table: %s", ftable.Name))
			} else {
				// delete the deleted rows
				err := targetDB.DeleteRows(ctx, tx, targetTable.Name, diffResult.DeletedRows, isDryRun)
				if err != nil {
					zap.L().Error("Error deleting rows ", zap.Error(err))
					return
				}

				zap.L().Info(fmt.Sprintf("Delete done on Table: %s", ftable.Name))
			}
		}

		if len(diffResult.UpdatedRows.ChangedColumns) > 0 {
			// update the updated rows
			err = targetDB.UpdateRows(ctx, tx, targetTable.Name, diffResult.UpdatedRows.ChangedColumns, diffResult.UpdatedRows.Before, diffResult.UpdatedRows.After, isDryRun)
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
			err := tx.Rollback()
			if err != nil {
				zap.S().Error("Failed to rollback transaction", zap.Error(err))
				return
			}
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
