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

type pr0nBot struct {
	bot reddit.Bot
}

func getSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

func ImageStreamFinder() {
	bot, _ := reddit.NewBotFromAgentFile("credentials.txt", 0)
	cfg := graw.Config{Subreddits: []string{"pics"}}
	handler := &pr0nBot{bot: bot}
	_, wait, _ := graw.Run(handler, bot, cfg)
	wait()
}

// Put the content to the ReadWriter's bucket at the key.
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

// Start initializes the bot.
func Start() {
	ImageStreamFinder()
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
