package database

import "database/sql"

func NewPostgresClient(connString string) (*sql.DB, error) {
	targetDB, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	return targetDB, err
}
