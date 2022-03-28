package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"stage-sync-cli/internal/propagate"
)

var scriptFile string
var dryRun bool
var logLevel string

func init() {
	runPropagation.PersistentFlags().StringVarP(&scriptFile, "script", "s", "", "path to script file")
	runPropagation.PersistentFlags().BoolVarP(&dryRun, "confirm", "c", false, "use to insert sql into database")
	runPropagation.PersistentFlags().StringVarP(&logLevel, "level", "l", "debug", "configure the minimal level of log output")
	runPropagation.MarkPersistentFlagRequired("script")
	viper.BindPFlag("confirm", runPropagation.PersistentFlags().Lookup("confirm"))
	viper.BindPFlag("script", runPropagation.PersistentFlags().Lookup("script"))
	viper.BindPFlag("level", runPropagation.PersistentFlags().Lookup("level"))

	rootCmd.AddCommand(runPropagation)
}

var runPropagation = &cobra.Command{
	Use:   "run",
	Short: "Runs the propagation",
	Run: func(cmd *cobra.Command, args []string) {
		var logger *zap.Logger
		if !dryRun && logLevel == "debug" {
			logger, _ = zap.NewDevelopment()
		} else {
			logger, _ = zap.NewProduction()
		}

		defer logger.Sync()
		zap.ReplaceGlobals(logger)

		if !dryRun {
			zap.L().Info("DRY RUN: not inserting into database")
		}
		propagate.Run(scriptFile, !dryRun)
	},
}
