package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theherk/pr0nbot/lib/scrape"
)

var startCmd = &cobra.Command{
	Use:           "start",
	Short:         "Start pr0nbot",
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		scrape.Start()
	},
}
