package builder

import (
	"context"
	"database/sql"
	"fmt"
	numeric "github.com/jackc/pgtype/ext/shopspring-numeric"
	"github.com/kilianstallz/stage-sync/internal/database"
	"github.com/kilianstallz/stage-sync/internal/database/postgres"
	"github.com/kilianstallz/stage-sync/internal/sql_log"
	"github.com/kilianstallz/stage-sync/models"
	"github.com/kilianstallz/stage-sync/pkg/config"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"reflect"
	"time"
)

type PostgresClient struct {
	connection *sql.DB
}

func (p *PostgresClient) NewConnection(credentials string) error {
	targetDB, err := sql.Open("postgres", credentials)
	if err != nil {
		panic(err)
	}
	p.connection = targetDB
	return err
}

func (p *PostgresClient) Close() error {
	err := p.connection.Close()
	return err
}

func (p *PostgresClient) Connection() *sql.DB {
	return p.connection
}

func (p *PostgresClient) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return p.connection.BeginTx(ctx, opts)
}

func (p *PostgresClient) BuildDeleteQuery(tableName string, rows models.Row) string {
	return postgres.BuildDeleteQuery(tableName, rows)
}

func (p *PostgresClient) BuildInsertQuery(tableName string, row models.Row) string {
	return postgres.BuildInsertQuery(tableName, row)
}

func (p *PostgresClient) BuildSelectQuery(table config.ConfigTable) string {
	return postgres.BuildSelectQuery(table)
}

func (p *PostgresClient) BuildUpdateQuery(tableName string, changedColumn []string, originalRow models.Row, updatedRow models.Row) string {
	return postgres.BuildUpdateQuery(tableName, changedColumn, originalRow, updatedRow)
}

func (p *PostgresClient) InsertRows(ctx context.Context, tx *sql.Tx, tableName string, rows []models.Row, isDryRun bool) error {
	_, f := sql_log.CreateSqlLogger("insert.sql")
	defer f.Close()

	for _, row := range rows {

		query := p.BuildInsertQuery(tableName, row)

		if !isDryRun {
			_, err := tx.ExecContext(ctx, query)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

func (p *PostgresClient) DeleteRows(ctx context.Context, tx *sql.Tx, tableName string, rows []models.Row, isDryRun bool) error {
	log, f := sql_log.CreateSqlLogger("delete.sql")
	defer f.Close()

	for _, row := range rows {

		query := p.BuildDeleteQuery(tableName, row)
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

func (p *PostgresClient) UpdateRows(ctx context.Context, tx *sql.Tx, tableName string, changedColumns []string, oldRows []models.Row, updatedRows []models.Row, isDryRun bool) error {
	log, f := sql_log.CreateSqlLogger("update.sql")
	defer f.Close()

	for i, oldRow := range oldRows {

		query := p.BuildUpdateQuery(tableName, changedColumns, oldRow, updatedRows[i])
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

func NewPostgresClient(credentials string) QueryBuilder {
	var client = PostgresClient{}
	err := client.NewConnection(credentials)
	if err != nil {
		panic(err)
	}
	return &client
}

func (p *PostgresClient) QueryTables(config *config.Config) []models.Table {
	var tables []models.Table
	for _, table := range config.Tables {
		tables = append(tables, p.buildTable(table))
	}

	return tables
}

func (p *PostgresClient) buildTable(table config.ConfigTable) models.Table {
	// Get all the data from the source database table
	q := p.BuildSelectQuery(table)
	// error using the interface implementation because of some type transformation in the background
	dbRows, err := p.Connection().Query(q)
	if err != nil {
		zap.S().Fatal("Error querying database", zap.Error(err))
	}
	defer dbRows.Close()

	// create table struct
	tableStruct := models.Table{
		Name:        table.Name,
		Rows:        []models.Row{},
		PrimaryKeys: table.PrimaryKeys,
		NoDelete:    table.NoDelete,
	}

	// iterate over the rows
	for dbRows.Next() {
		// rows := make([]models.Column, 0, len(columns))
		cols, _ := dbRows.ColumnTypes()
		pointerValues := make([]interface{}, len(cols))
		object := map[string]interface{}{}
		for i, column := range cols {
			typeName := column.DatabaseTypeName()
			isNullable := true
			switch typeName {
			case "INT":
				if isNullable {
					pointerValues[i] = new(sql.NullInt64)
					object[column.Name()] = new(sql.NullInt64)
				} else {
					pointerValues[i] = new(int)
					object[column.Name()] = new(int)
				}
			case "VARCHAR":
				if isNullable {
					pointerValues[i] = new(sql.NullString)
					object[column.Name()] = new(sql.NullString)
				} else {
					pointerValues[i] = new(string)
					object[column.Name()] = new(string)
				}
			case "DATETIME":
				if isNullable {
					pointerValues[i] = new(sql.NullTime)
					object[column.Name()] = new(sql.NullTime)
				} else {
					pointerValues[i] = new(time.Time)
					object[column.Name()] = new(time.Time)
				}
			case "NUMERIC":
				pointerValues[i] = new(numeric.Numeric)
				object[column.Name()] = new(numeric.Numeric)
			case "DECIMAL":
				pointerValues[i] = new(decimal.NullDecimal)
				object[column.Name()] = new(decimal.NullDecimal)
			case "TEXT":
				if isNullable {
					pointerValues[i] = new(sql.NullString)
					object[column.Name()] = new(sql.NullString)
				} else {
					pointerValues[i] = new(string)
					object[column.Name()] = new(string)
				}
			case "BOOL":
				if isNullable {
					pointerValues[i] = new(sql.NullBool)
					object[column.Name()] = new(sql.NullBool)
				} else {
					pointerValues[i] = new(bool)
					object[column.Name()] = new(bool)
				}
			case "DATE":
				if isNullable {
					pointerValues[i] = new(sql.NullTime)
					object[column.Name()] = new(sql.NullTime)
				} else {
					pointerValues[i] = new(time.Time)
					object[column.Name()] = new(time.Time)
				}
			case "INT4":
				if isNullable {
					pointerValues[i] = new(sql.NullInt64)
					object[column.Name()] = new(sql.NullInt64)
				} else {
					pointerValues[i] = new(int)
					object[column.Name()] = new(int)
				}

			case "TIMESTAMP":
				if isNullable {
					pointerValues[i] = new(sql.NullTime)
					object[column.Name()] = new(sql.NullTime)
				} else {
					pointerValues[i] = new(time.Time)
					object[column.Name()] = new(time.Time)
				}
			default:
				panic(fmt.Sprintf("Unsupported type: %s", typeName))
			}
		}
		err := dbRows.Scan(pointerValues...)
		if err != nil {
			panic(err)
		}
		// pointer values is an array of pointers to the values
		// make an array of actual values
		values := make([]interface{}, len(pointerValues))
		for i, pointerValue := range pointerValues {
			values[i] = reflect.ValueOf(pointerValue).Elem().Interface()
		}

		row := make([]models.Column, 0, len(table.Columns))
		for i, column := range table.Columns {
			val := database.ConvertDbValue(values[i])
			typ := reflectType(val)
			row = append(row, models.Column{
				Name:  column,
				Value: val,
				Type:  typ,
			})
		}
		tableStruct.Rows = append(tableStruct.Rows, row)
	}

	return tableStruct
}

func reflectType(v interface{}) string {
	if v == nil {
		return "NULL"
	}
	return reflect.TypeOf(v).String()
}
