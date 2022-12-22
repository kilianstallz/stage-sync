package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"stage-sync/internal/cli/base"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	cmd := &cobra.Command{Use: "stage-sync"}

	cmd.AddCommand(base.InitCmd(), base.VersionCmd())

	_ = cmd.Execute()
}
