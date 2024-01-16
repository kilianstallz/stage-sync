package models

// Table represents a PostgresSQL table
type Table struct {
	Name        string
	Rows        []Row
	PrimaryKeys []string
	NoDelete    bool
}

type Row []Column

func (r Row) GetColumn(name string) *Column {
	for _, c := range r {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

// Column is an in memory representation of a column.
type Column struct {
	Name  string
	Type  string
	Value interface{}
	IsPK  bool
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
