package base

import (
	"fmt"
	"github.com/kilianstallz/stage-sync/internal"
	"github.com/kilianstallz/stage-sync/pkg/propagation"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func InitCmd() *cobra.Command {
	var scriptFilePath, logLevel string
	var execute bool
	var source string
	var target string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs the propagation",
		Run: func(cmd *cobra.Command, args []string) {
			var logger *zap.Logger
			logger, _ = zap.NewDevelopment()
			defer logger.Sync()
			zap.ReplaceGlobals(logger)

			if !execute {
				zap.L().Info("DRY RUN: not inserting into database")
			}

			err := propagation.Execute(scriptFilePath, execute, source, target)
			cobra.CheckErr(err)
		},
	}

	cmd.PersistentFlags().StringVarP(&scriptFilePath, "script", "s", "", "path to script file")
	cmd.PersistentFlags().StringVarP(&logLevel, "level", "l", "debug", "configure the minimal level of log output")
	cmd.PersistentFlags().BoolVarP(&execute, "confirm", "c", false, "Execute propagation without dry run")
	cmd.PersistentFlags().StringVar(&source, "source", "", "The source database stage")
	cmd.PersistentFlags().StringVar(&target, "target", "", "The target database stage")
	cmd.MarkPersistentFlagRequired("script")

	return cmd
}

func VersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use: "version",

		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(internal.Version)
		},
	}

	return versionCmd
}
