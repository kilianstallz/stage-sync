package database

import "database/sql"

func NewPostgresClient(connString string) (*sql.DB, error) {
	targetDB, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	return targetDB, err
}

func NewSqliteClient() (*sql.DB, error) {
	targetDB, err := sql.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}
	return targetDB, err
}
