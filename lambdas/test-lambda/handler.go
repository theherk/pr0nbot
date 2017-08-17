package main

import (
	"errors"

	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

// Event is a downtoearth event.
type Event struct {
	Body        map[string]interface{} `json:"body"`
	Path        map[string]string      `json:"path"`
	Querystring map[string]string      `json:"querystring"`
	Route       string                 `json:"route"`
}

// Response is an object with the data to return.
type Response struct {
	Messages []string
}

// Handle lambda event.
func Handle(evt Event, ctx *runtime.Context) (Response, error) {
	return Response{}, errors.New("no matching route")
}
