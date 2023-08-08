package sql_log

import (
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
)

func CreateSqlLogger(path string) (*log.Logger, *os.File) {
	if path == "" {
		path = "./sql.log"
	}
	f, err := os.OpenFile(filepath.Clean(path), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		zap.S().Fatalf("error opening file: %v", err)
	}
	return log.New(f, "", 0), f
}
