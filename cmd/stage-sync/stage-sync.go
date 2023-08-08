package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kilianstallz/stage-sync/internal"
	"github.com/kilianstallz/stage-sync/internal/cli/base"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	cmd := &cobra.Command{
		Use:     "stage-sync",
		Version: fmt.Sprint(internal.Version),
	}

	cmd.AddCommand(base.InitCmd())

	_ = cmd.Execute()
}
