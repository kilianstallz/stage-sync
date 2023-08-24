package adapter

import "database/sql"

type DatabaseConnector interface {
	NewConnection(connString string) (*sql.DB, error)
	Disconnect() error
}
