package table

import (
	"reflect"
	"stage-sync/models"
)

func FindTable(tables []models.Table, name string) models.Table {
	for _, table := range tables {
		if table.Name == name {
			return table
		}
	}
	return models.Table{}
}

func FindColumn(row models.Row, name string) *models.Column {
	for _, column := range row {
		if column.Name == name {
			return &column
		}
	}
	return nil
}

// FindRow finds and returns a row where all the primary keys match the values in the row.
func FindRow(rows []models.Row, keys []string, row models.Row) *models.Row {
	findPkColsMap := make(map[string]interface{})
	for _, key := range keys {
		col := *FindColumn(row, key)
		findPkColsMap[key] = col.Value
	}

	for _, row := range rows {

		searchPKColsMap := make(map[string]interface{})
		for _, key := range keys {
			col := *FindColumn(row, key)
			searchPKColsMap[key] = col.Value
		}

		keyArray := make([]bool, len(keys))
		for i, key := range keys {
			if reflect.DeepEqual(findPkColsMap[key], searchPKColsMap[key]) {
				keyArray[i] = true
			}
		}

		if allTrue(keyArray) {
			return &row
		}

	}
	return nil
}

func allTrue(array []bool) bool {
	for _, value := range array {
		if !value {
			return false
		}
	}
	return true
}
