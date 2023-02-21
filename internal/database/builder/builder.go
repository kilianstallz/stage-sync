package builder

import (
	"context"
	"database/sql"
	"github.com/kilianstallz/stage-sync/models"
	"github.com/kilianstallz/stage-sync/pkg/config"
)

type QueryBuilder interface {
	NewConnection(credentials string) error
	Close() error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Connection() *sql.DB

	BuildDeleteQuery(tableName string, rows models.Row) string
	BuildInsertQuery(tableName string, row models.Row) string
	BuildSelectQuery(table config.ConfigTable) string
	BuildUpdateQuery(tableName string, changedColumn []string, originalRow models.Row, updatedRow models.Row) string

	InsertRows(ctx context.Context, tx *sql.Tx, tableName string, rows []models.Row, isDryRun bool) error
	UpdateRows(ctx context.Context, tx *sql.Tx, tableName string, changedColumns []string, oldRows []models.Row, updatedRows []models.Row, isDryRun bool) error
	DeleteRows(ctx context.Context, tx *sql.Tx, tableName string, rows []models.Row, isDryRun bool) error
	QueryTables(config *config.Config) (tables []models.Table, err error)
}
