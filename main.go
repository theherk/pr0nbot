package main

import (
	"fmt"
	"os"

	"github.com/theherk/pr0nbot/cmd"
)

type announcer struct{}

func main() {
	if err := cmd.Pr0nbotCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
