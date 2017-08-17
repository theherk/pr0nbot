package cmd

import "github.com/spf13/cobra"

// RootCmd is the main command.
var RootCmd = &cobra.Command{
	Use:   "pr0nbot",
	Short: "pr0nbot short description",
	Long:  "pr0nbot long description.",
}

func init() {
	RootCmd.AddCommand(sampleCmd, versionCmd)
}
