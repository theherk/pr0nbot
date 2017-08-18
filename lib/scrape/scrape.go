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
	"github.com/theherk/pr0nbot/lib/comment"
	"github.com/theherk/pr0nbot/lib/detect"
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
	_, wait, err := graw.Run(handler, bot, cfg)
	if err != nil {
		fmt.Println(err)
	}
	if err := wait(); err != nil {
		fmt.Println(err)
	}
}

// Put the content to the ReadWriter's bucket at the key.
func Put(content io.Reader, length int64, key string) error {
	svc := s3.New(getSession())
	input := &s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(content),
		Bucket:               aws.String(bucket),
		ContentLength:        aws.Int64(length),
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
		resp, err := http.Get(p.URL)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(resp)
		// send a put to S3 bucket
		if err := Put(resp.Body, resp.ContentLength, p.Name); err != nil {
			return err
		}
		is, err := detect.IsPrawn(p.URL)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if is {
			comment.Do(r.bot, p)
		} else {
			fmt.Println(":( not prawn")
		}
	}
	return nil
}
