// Package scrape gets new posts from given subreddits and passes them to detect.
package scrape

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

const bucket = "prawnbot"

// func ImageFinder() (string, error) {
// 	bot, err := reddit.NewBotFromAgentFile("credentials.txt", 0)
// 	harvest, err := bot.Listing("/r/pics", "")
// 	if err != nil {
// 		fmt.Println("Failed to fetch /r/golang: ", err)
// 		return "", err
// 	}

// 	for _, post := range harvest.Posts[:5] {
// 		if strings.HasSuffix(post.URL, ".jpg") {
// 			// read file from URL in memory
// 			resp, e := http.Get(post.URL)
// 			// send a put to S3 bucket
// 			Put(resp.Body, post.URL)
// 			fmt.Printf("[%s] posted [%s]\n", post.Title, post.URL)
// 		}
// 	}
// 	return post.
// }

type pr0nBot struct {
	bot reddit.Bot
}

func ImageStreamFinder() {
	bot, _ := reddit.NewBotFromAgentFile("credentials.txt", 0)
	cfg := graw.Config{Subreddits: []string{"pics"}}
	handler := &pr0nBot{bot: bot}
	_, wait, _ := graw.Run(handler, bot, cfg)
	wait()
}

func (r *pr0nBot) Post(p *reddit.Post) error {
	if strings.HasSuffix(p.URL, ".jpg") {
		<-time.After(10 * time.Second)
		fmt.Printf("Image: %s\n", p.URL)
		// read file from URL in memory
		resp, _ := http.Get(p.URL)
		// send a put to S3 bucket
		Put(resp.Body, p.URL)
	}
	return nil
}

//Put puts the content to the ReadWriter's bucket at the key
func Put(content io.Reader, key string) error {
	svc := s3.New(getSession())
	input := &s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(content),
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		ServerSideEncryption: aws.String("AES256"),
	}
	if _, err := svc.PutObject(input); err != nil {
		return fmt.Errorf("experienced error putting %s", err)
	}
	return nil
}

func getSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}
