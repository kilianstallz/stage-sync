package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"stage-sync/internal"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(internal.Version)
	},
}
