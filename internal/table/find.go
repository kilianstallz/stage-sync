package table

import (
	"github.com/kilianstallz/stage-sync/pkg/models"
	"github.com/samber/lo"
	"reflect"
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
	c, ok := lo.Find[models.Column](row, func(column models.Column) bool {
		return column.Name == name
	})
	if !ok {
		return nil
	}
	return &c
}

// FindRow finds and returns a row where all the primary keys match the values in the row.
func FindRow(rows []models.Row, keys []string, row models.Row) *models.Row {
	findPkColsMap := make(map[string]interface{})
	for _, key := range keys {
		col := FindColumn(row, key)
		if col != nil {
			findPkColsMap[key] = col.Value
		}
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

		if AllTrue(keyArray) {
			return &row
		}

	}
	return nil
}

func AllTrue(array []bool) bool {
	b := lo.IndexOf[bool](array, false)
	if b > -1 { // -1 not found
		return false
	}
	return true
}
