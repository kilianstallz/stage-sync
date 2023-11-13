package propagation

import (
	"context"
	"fmt"
	"github.com/kilianstallz/stage-sync/internal/database/builder"
	"github.com/kilianstallz/stage-sync/internal/diff"
	"github.com/kilianstallz/stage-sync/internal/table"
	"github.com/kilianstallz/stage-sync/pkg/config"
	"go.uber.org/zap"
)

func Execute(configPath string, execute bool, source, target string) error {
	zap.L().Info("Starting propagation")
	conf, _ := config.ParseConfigFromFile(configPath)

	sourceDatabase, ok := conf.Stages[source]
	targetDatabase, tok := conf.Stages[target]

	if !ok {
		zap.S().Fatal("Source stage not found in config")
	}

	if !tok {
		zap.S().Fatal("Target stage not found in config")
	}

	// Get the source database connection.
	sourceDB := builder.NewPostgresClient(sourceDatabase)
	defer sourceDB.Close()

	tables, err := sourceDB.QueryTables(conf)
	if err != nil {
		return err
	}

	err = sourceDB.Close()
	if err != nil {
		zap.S().Error("Failed to close source database connection", zap.Error(err))
	}

	// Get the target database connection.
	targetDB := builder.NewPostgresClient(targetDatabase)
	defer targetDB.Close()

	targetTables, err := targetDB.QueryTables(conf)
	if err != nil {
		return err
	}

	defer func(targetDB builder.QueryBuilder) {
		err := targetDB.Close()
		if err != nil {
			zap.S().Error("Failed to close target database connection", zap.Error(err))
		}
	}(targetDB)

	ctx := context.Background()

	tx, err := targetDB.BeginTx(ctx, nil)
	if err != nil {
		return nil
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
			err := targetDB.InsertRows(ctx, tx, targetTable.Name, diffResult.AddedRows, !execute)
			if err != nil {
				zap.L().Error("Error inserting rows ", zap.Error(err))
				return err
			}
			zap.L().Debug(fmt.Sprintf("Insert done on Table: %s", ftable.Name))

		}

		if len(diffResult.DeletedRows) > 0 {
			if ftable.NoDelete {
				zap.L().Debug(fmt.Sprintf("No delete on Table: %s", ftable.Name))
			} else {
				// delete the deleted rows
				err := targetDB.DeleteRows(ctx, tx, targetTable.Name, diffResult.DeletedRows, !execute)
				if err != nil {
					zap.L().Error("Error deleting rows ", zap.Error(err))
					return err
				}

				zap.L().Debug(fmt.Sprintf("Delete done on Table: %s", ftable.Name))
			}
		}

		if len(diffResult.UpdatedRows.ChangedColumns) > 0 {
			// update the updated rows
			err = targetDB.UpdateRows(ctx, tx, targetTable.Name, diffResult.UpdatedRows.ChangedColumns, diffResult.UpdatedRows.Before, diffResult.UpdatedRows.After, !execute)
			if err != nil {
				zap.S().Error("Error updating rows ", zap.Error(err))
				return err
			}
			zap.L().Info(fmt.Sprintf("Update done on Table: %s", ftable.Name))
		}
	}

	if execute {
		// commit the transaction
		err = tx.Commit()
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				zap.S().Error("Failed to rollback transaction", zap.Error(err))
				return err
			}
			zap.S().Error("Failed to commit transaction", zap.Error(err))
			return err
		}
	} else {
		zap.L().Info("Dry run, skipping transaction commit")
		err = tx.Rollback()
		if err != nil {
			zap.S().Error("Failed to rollback transaction", zap.Error(err))
			return err
		}
	}
	return nil
}
