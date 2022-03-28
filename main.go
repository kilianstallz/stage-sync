package main

import (
	"go.uber.org/zap"
	"stage-sync-cli/cmd"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	cmd.Execute()
}
