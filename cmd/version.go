package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theherk/pr0nbot/util"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current pr0nbot version.",
	Long: "This will be a semantic version with the version of source " +
		"from which this application was built.",
	Run: func(cmd *cobra.Command, args []string) {
		util.ShowVersion()
	},
}
