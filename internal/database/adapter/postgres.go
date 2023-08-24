package adapter

import "github.com/kilianstallz/stage-sync/internal/database/builder"

type PostgresAdapter struct {
	builder.QueryBuilder
	DatabaseConnector
}

func NewPostgresAdapter() *PostgresAdapter {
	return &PostgresAdapter{}
}
