package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theherk/pr0nbot/lib"
)

var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Run the sample command.",
	Long:  "Running sample shows an example cli implementation.",
	Run: func(cmd *cobra.Command, args []string) {
		lib.Sample()
	},
}
