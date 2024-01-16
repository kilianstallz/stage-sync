package builder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"reflect"
	"time"

	numeric "github.com/jackc/pgtype/ext/shopspring-numeric"
	"github.com/kilianstallz/stage-sync/internal/database"
	"github.com/kilianstallz/stage-sync/internal/database/postgres"
	"github.com/kilianstallz/stage-sync/internal/sql_log"
	"github.com/kilianstallz/stage-sync/pkg/config"
	"github.com/kilianstallz/stage-sync/pkg/models"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type PostgresClient struct {
	connection *sql.DB
}

func (p *PostgresClient) NewConnection(credentials string) error {
	targetDB, err := sql.Open("pgx", credentials)
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
	log, f := sql_log.CreateSqlLogger("insert.sql")
	defer f.Close()

	for _, row := range rows {

		query := p.BuildInsertQuery(tableName, row)
		log.Println(query)

		if !isDryRun {
			_, err := tx.ExecContext(ctx, query)
			if err != nil {
				// check if the error string contains the error code for foreign key constraint violation
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					if pgErr.Code == "23503" {
						zap.S().Debug("Foreign key constraint violation")
						//zap.S().Debug(pgErr.ColumnName)
						//// find col where name is ParentId
						//c := row.GetColumn("ParentId")
						//// parse value as string
						//parentId := c.Value.(string)
						//zap.S().Debug(parentId)
						//// find row where PK is parentId
						//for _, r := range rows {
						//	if r.GetColumn("Id").Value.(string) == parentId {
						//		q := p.BuildInsertQuery(tableName, r)
						//		zap.S().Debug(q)
						//	}
						//}

						// TODO: Schedule the row for later insertion
					}
				}

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
		zap.L().Debug("Error connecting to database", zap.Error(err))
		log.Fatalf("Error connecting to database: %v", err)
	}
	return &client
}

func (p *PostgresClient) QueryTables(config *config.Config) ([]models.Table, error) {
	var tables []models.Table
	for _, table := range config.Tables {
		tables = append(tables, p.BuildTable(table))
	}

	return tables, nil
}

func (p *PostgresClient) BuildTable(table config.ConfigTable) models.Table {
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
		cols, _ := dbRows.ColumnTypes()
		pointerValues := make([]interface{}, len(cols))
		for i, column := range cols {
			typeName := column.DatabaseTypeName()
			isNullable := true
			if handler, ok := typeHandlers[typeName]; ok {
				// handler returns a pointer to the actual value
				// if the type is not nullable, the pointer is not needed
				// so we dereference it
				//if typeName == "_INT4" {
				//	pointerValues[i] = reflect.ValueOf(handler(isNullable)).Elem().Interface()
				//} else {
				pointerValues[i] = handler(isNullable)
				//}
			} else {
				panic(fmt.Sprintf("No handler for type %s", typeName))
			}
		}
		err := dbRows.Scan(pointerValues...)
		if err != nil {
			panic(err)
		}
		// pointer values is an array of array of actual valuepointers to the values
		// make an s
		values := make([]interface{}, len(pointerValues))
		for i, pointerValue := range pointerValues {
			values[i] = reflect.ValueOf(pointerValue).Elem().Interface()
		}

		row := make([]models.Column, 0, len(table.Columns))
		for i, column := range table.Columns {
			val := database.ConvertDbValue(values[i])
			typ := reflectType(val)
			isPK := isInArray(column, table.PrimaryKeys)
			row = append(row, models.Column{
				Name:  column,
				Value: val,
				Type:  typ,
				IsPK:  isPK,
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

var typeHandlers = map[string]func(isNullable bool) interface{}{
	"INT": func(isNullable bool) interface{} {
		if isNullable {
			return new(sql.NullInt64)
		} else {
			return new(int)
		}
	},
	"VARCHAR": func(isNullable bool) interface{} {
		if isNullable {
			return new(sql.NullString)
		} else {
			return new(string)
		}
	},
	"DATETIME": func(isNullable bool) interface{} {
		if isNullable {
			return new(sql.NullTime)
		} else {
			return new(time.Time)
		}
	},
	"NUMERIC": func(isNullable bool) interface{} {
		return new(numeric.Numeric)
	},
	"DECIMAL": func(isNullable bool) interface{} {
		if isNullable {
			return new(decimal.NullDecimal)
		} else {
			return new(decimal.Decimal)
		}
	},
	"TEXT": func(isNullable bool) interface{} {
		if isNullable {
			return new(sql.NullString)
		} else {
			return new(string)
		}
	},
	"BOOL": func(isNullable bool) interface{} {
		if isNullable {
			return new(sql.NullBool)
		} else {
			return new(bool)
		}
	},
	"DATE": func(isNullable bool) interface{} {
		if isNullable {
			return new(sql.NullTime)
		} else {
			return new(time.Time)
		}
	},
	"INT4": func(isNullable bool) interface{} {
		if isNullable {
			return new(sql.NullInt64)
		} else {
			return new(int)
		}
	},
	"TIMESTAMP": func(isNullable bool) interface{} {
		if isNullable {
			return new(sql.NullTime)
		} else {
			return new(time.Time)
		}
	},
	// handle _INT4 arrays as []int
	"_INT4": func(isNullable bool) interface{} {
		return new(pgtype.Int4Array)
	},
}

func isInArray(target string, array []string) bool {
	for _, element := range array {
		if element == target {
			return true
		}
	}
	return false
}
