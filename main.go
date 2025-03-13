package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request: method=%s, path=%s\n", req.HTTPMethod, req.Path)
	fmt.Printf("Headers: %v\n", req.Headers)

	path := req.Path
	if path == "" {
		path = "/"
	}

	if req.HTTPMethod == "GET" && path == "/hello" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"message":"Hello World from the Dog Walking API!"}`,
		}, nil
	}

	if req.HTTPMethod == "POST" && path == "/dogs" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{"error":"Route not found","path":"%s","method":"%s"}`, path, req.HTTPMethod),
	}, nil
}

func main() {
	lambda.Start(handler)
}
