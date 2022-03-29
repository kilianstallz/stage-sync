package models

// Table represents a PostgresSQL table
type Table struct {
	Name        string
	Rows        []Row
	PrimaryKeys []string
	NoDelete bool
}

type Row []Column

// Column is an in memory representation of a column.
type Column struct {
	Name  string
	Type  string
	Value interface{}
}

type DiffResult struct {
	AddedRows   []Row
	DeletedRows []Row
	UpdatedRows struct {
		Before         []Row
		After          []Row
		ChangedColumns []string
	}
}
