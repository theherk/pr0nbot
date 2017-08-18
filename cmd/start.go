package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/theherk/pr0nbot/lib/scrape"
)

var startCmd = &cobra.Command{
	Use:           "start [subs]",
	Short:         "Start pr0nbot",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("start requires at least one sub")
		}
		scrape.Start(args)
		return nil
	},
}
