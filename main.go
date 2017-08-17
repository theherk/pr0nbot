package main

import (
	"fmt"
	"strings"

	"github.com/turnage/graw/reddit"
)

type announcer struct{}

func main() {
	bot, err := reddit.NewBotFromAgentFile("credentials.txt", 0)
	harvest, err := bot.Listing("/r/pics", "")
	if err != nil {
		fmt.Println("Failed to fetch /r/golang: ", err)
		return
	}

	for _, post := range harvest.Posts[:5] {
		if strings.HasSuffix(post.URL, ".jpg") {
			fmt.Printf("[%s] posted [%s]\n", post.Title, post.URL)
		}
	}
	// if err := cmd.RootCmd.Execute(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}
