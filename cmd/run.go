package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"stage-sync-cli/internal/propagate"
)

var scriptFile string
var dryRun bool

func init() {
	runPropagation.PersistentFlags().StringVarP(&scriptFile, "script", "s", "", "path to script file")
	runPropagation.PersistentFlags().BoolVarP(&dryRun, "confirm", "c", true, "use to insert sql into database")
	runPropagation.MarkPersistentFlagRequired("script")
	viper.BindPFlag("confirm", runPropagation.PersistentFlags().Lookup("confirm"))
	viper.BindPFlag("script", runPropagation.PersistentFlags().Lookup("script"))

	rootCmd.AddCommand(runPropagation)
}

var runPropagation = &cobra.Command{
	Use:   "run",
	Short: "Runs the propagation",
	Run: func(cmd *cobra.Command, args []string) {
		if !dryRun {
			zap.L().Info("Use confirm argument to insert sql into database")
		}
		propagate.Run(scriptFile, !dryRun)
	},
}
