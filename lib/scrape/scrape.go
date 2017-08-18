// Package scrape gets new posts from given subreddits and passes them to detect.
package scrape

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

type pr0nBot struct {
	cfg  graw.Config
	rbot reddit.Bot
	wait func() error
}

const bucket = "prawnbot"

func getSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

// Put the content to S3.
func Put(content io.Reader, length int64, key string) error {
	svc := s3.New(getSession())
	fmt.Println(key)
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
func Start(subs []string) {
	rbot, err := reddit.NewBotFromAgentFile("credentials.txt", 0)
	if err != nil {
		panic(err)
	}
	bot := &pr0nBot{
		cfg:  graw.Config{Subreddits: subs},
		rbot: rbot,
	}
	if err := bot.initWait(); err != nil {
		panic(err)
	}
	bot.run()
}

func (p *pr0nBot) initWait() error {
	_, wait, err := graw.Run(p, p.rbot, p.cfg)
	if err != nil {
		return err
	}
	p.wait = wait
	return nil
}

func (r *pr0nBot) Post(p *reddit.Post) error {
	if strings.HasSuffix(p.URL, ".jpg") {
		<-time.After(10 * time.Second)
		fmt.Printf("Image: %s\n", p.URL)
		resp, err := http.Get(p.URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		reader := bytes.NewReader(body)
		if err := Put(reader, resp.ContentLength, p.Name); err != nil {
			return err
		}
		is, err := detect.IsPrawn(p.Name)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if is {
			comment.Do(r.rbot, p)
		} else {
			fmt.Println(":( not prawn")
		}
	}
	return nil
}

func (p *pr0nBot) run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered err in wait: ", r)
		}
	}()
	if err := p.wait(); err != nil {
		panic(err)
	}
}
