package cmd

import "github.com/spf13/cobra"

// RootCmd is the main command.
var Pr0nbotCmd = &cobra.Command{
	Use:   "pr0nbot",
	Short: "pr0nbot detects prawn images and thinks of the children",
}

func init() {
	Pr0nbotCmd.AddCommand(startCmd, versionCmd)
}
