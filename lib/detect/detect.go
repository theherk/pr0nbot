// Package detect determines if the found image contains prawn.
package detect

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

const Bucket = "prawnbot"

func getLables(bucket, key string) ([]*rekognition.Label, error) {
	svc := rekognition.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))
	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(key),
			},
		},
		MaxLabels:     aws.Int64(100),
		MinConfidence: aws.Float64(50.000000),
	}
	out, err := svc.DetectLabels(input)
	if err != nil {
		return nil, err
	}
	return out.Labels, nil
}

func IsPrawn(key string) (bool, error) {
	res := false
	labels, err := getLables(Bucket, key)
	if err != nil {
		return res, err
	}
	if contains("shrimp", labels) {
		res = true
	}
	if contains("crab", labels) {
		res = true
	}
	if contains("prawn", labels) {
		res = true
	}
	if contains("lobster", labels) {
		res = true
	}
	if contains("insect", labels) {
		if contains("fish", labels) {
			res = true
		}
		if contains("ocean", labels) {
			res = true
		}
		if contains("sea", labels) {
			res = true
		}
		if contains("sea life", labels) {
			res = true
		}
	}
	if contains("lice", labels) {
		if contains("fish", labels) {
			res = true
		}
		if contains("ocean", labels) {
			res = true
		}
		if contains("sea", labels) {
			res = true
		}
		if contains("sea life", labels) {
			res = true
		}
	}
	return res, nil
}

func contains(word string, labels []*rekognition.Label) bool {
	contains := false
	for _, l := range labels {
		if word == strings.ToLower(*l.Name) {
			contains = true
		}
	}
	return contains
}
