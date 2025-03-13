package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

type Dog struct {
	Name string `json:"name"`
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	fmt.Printf("Request: method=%s, path=%s, rawPath=%s\n",
		req.RequestContext.HTTP.Method,
		req.RequestContext.HTTP.Path,
		req.RawPath)
	fmt.Printf("Headers: %v\n", req.Headers)
	fmt.Printf("Route: %s\n", req.RouteKey)

	path := req.RequestContext.HTTP.Path
	method := req.RequestContext.HTTP.Method

	if method == "GET" && path == "/ping" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
			Body: "OK",
		}, nil
	}

	if method == "POST" && path == "/dogs" {
		var dog Dog
		err := json.Unmarshal([]byte(req.Body), &dog)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf(`{"error":"%s"}`, err.Error()),
			}, nil
		}

		returnedDog := Dog{Name: dog.Name}
		returnedDogJson, err := json.Marshal(returnedDog)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf(`{"error":"%s"}`, err.Error()),
			}, nil
		}

		return events.APIGatewayV2HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(returnedDogJson),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 404,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{"error":"Route not found","path":"%s","method":"%s","routeKey":"%s"}`,
			path, method, req.RouteKey),
	}, nil
}

func main() {
	lambda.Start(handler)
}
