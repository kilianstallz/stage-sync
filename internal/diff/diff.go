package diff

import (
	"fmt"

	"github.com/kilianstallz/stage-sync/internal/table"
	"github.com/kilianstallz/stage-sync/pkg/models"
	"go.uber.org/zap"
)

func FindDiffResult(source models.Table, target models.Table) models.DiffResult {
	primaryKeys := source.PrimaryKeys

	sourceTable := source
	targetTable := target

	diffResult := &models.DiffResult{
		AddedRows:   make([]models.Row, 0),
		DeletedRows: make([]models.Row, 0),
		UpdatedRows: struct {
			Before         []models.Row
			After          []models.Row
			ChangedColumns []string
		}{Before: make([]models.Row, 0), After: make([]models.Row, 0), ChangedColumns: []string{}},
	}

	// Diff the added rows based on the primary keys using the FindRow function
	zap.L().Info("Calculating diff on Table: " + sourceTable.Name)
	FindAddedAndChangedRows(sourceTable, targetTable, primaryKeys, diffResult)

	// Diff the deleted rows based on the primary keys using the FindRow function
	FindRemovedRows(targetTable, sourceTable, primaryKeys, diffResult)

	return *diffResult
}

func FindRemovedRows(targetTable models.Table, sourceTable models.Table, primaryKeys []string, diffResult *models.DiffResult) {
	for _, targetRow := range targetTable.Rows {
		sourceRow := table.FindRow(sourceTable.Rows, primaryKeys, targetRow)
		if sourceRow == nil {
			diffResult.DeletedRows = append(diffResult.DeletedRows, targetRow)
		}
	}
}

func FindAddedAndChangedRows(sourceTable models.Table, targetTable models.Table, primaryKeys []string, diffResult *models.DiffResult) {
	targetRowsMap := make(map[string]models.Row)
	for _, row := range targetTable.Rows {
		key := generateKeyFromRow(row, primaryKeys)
		targetRowsMap[key] = row
	}

	for _, sourceRow := range sourceTable.Rows {
		key := generateKeyFromRow(sourceRow, primaryKeys)
		targetRow, found := targetRowsMap[key]
		if !found {
			diffResult.AddedRows = append(diffResult.AddedRows, sourceRow)
		} else {
			for _, sourceColumn := range sourceRow {
				targetColumn := table.FindColumn(targetRow, sourceColumn.Name)
				if targetColumn != nil && sourceColumn.Value != targetColumn.Value {
					diffResult.UpdatedRows.ChangedColumns = append(diffResult.UpdatedRows.ChangedColumns, targetColumn.Name)
					diffResult.UpdatedRows.Before = append(diffResult.UpdatedRows.Before, targetRow)
					diffResult.UpdatedRows.After = append(diffResult.UpdatedRows.After, sourceRow)
					break
				}
			}
		}
	}
	// lo.ForEach[models.Row](sourceTable.Rows, func(sourceRow models.Row, _ int) {

	// 	targetRow := table.FindRow(targetTable.Rows, primaryKeys, sourceRow)

	// 	if targetRow == nil {
	// 		diffResult.AddedRows = append(diffResult.AddedRows, sourceRow)
	// 	} else {
	// 		lo.ForEach[models.Column](sourceRow, func(sourceColumn models.Column, _ int) {
	// 			targetColumn := table.FindColumn(*targetRow, sourceColumn.Name)
	// 			if targetColumn != nil {
	// 				if sourceColumn.Value != targetColumn.Value {
	// 					diffResult.UpdatedRows.ChangedColumns = append(diffResult.UpdatedRows.ChangedColumns, targetColumn.Name)
	// 					diffResult.UpdatedRows.Before = append(diffResult.UpdatedRows.Before, *targetRow)
	// 					diffResult.UpdatedRows.After = append(diffResult.UpdatedRows.After, sourceRow)
	// 				}
	// 			}
	// 		})
	// 	}
	// })
}

func generateKeyFromRow(row models.Row, primaryKeys []string) string {
	var key string
	for _, column := range row {
		if contains(primaryKeys, column.Name) {
			key += fmt.Sprintf("%v", column.Value)
		}
	}
	return key
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
