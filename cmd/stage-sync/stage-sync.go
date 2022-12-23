package main

import (
	"github.com/kilianstallz/stage-sync/internal/cli/base"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	cmd := &cobra.Command{Use: "stage-sync"}

	cmd.AddCommand(base.InitCmd(), base.VersionCmd())

	_ = cmd.Execute()
}
