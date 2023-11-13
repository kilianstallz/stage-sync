package builder

import (
	"context"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kilianstallz/stage-sync/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresClient_QueryTables(t *testing.T) {

	ctx := context.Background()

	// Create db using testcontainers
	pgContainer, err := postgres.RunContainer(ctx, testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithDatabase("source"), testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)))

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	client := NewPostgresClient(connStr)
	defer func(client QueryBuilder) {
		err := client.Close()
		if err != nil {
			t.Error(err)
		}
	}(client)

	// Create a table in the database
	_, err = client.Connection().Exec("CREATE TABLE test (id INT PRIMARY KEY, name VARCHAR(50))")
	assert.NoError(t, err)

	// Insert some data into the table
	_, err = client.Connection().Exec("INSERT INTO test (id, name) VALUES (1, 'John'), (2, 'Jane')")
	assert.NoError(t, err)

	// Call buildTable
	table := config.ConfigTable{
		Name:        "test",
		PrimaryKeys: []string{"id"},
		NoDelete:    false,
		Columns:     []string{"id", "name"},
	}
	result := client.BuildTable(table)

	// print the result
	t.Log(result.Name)

	// Check the result
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 2, len(result.Rows))
	assert.Equal(t, []string{"id"}, result.PrimaryKeys)
	assert.False(t, result.NoDelete)
}
