package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Dog struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	path := req.RequestContext.HTTP.Path
	method := req.RequestContext.HTTP.Method
	fmt.Printf("Request: %s %s\n", method, path)

	if method == "GET" && path == "/ping" {
		return healthCheck()
	}

	if method == "GET" && path == "/dogs" {
		return listDogs()
	}

	if method == "POST" && path == "/dogs" {
		return postDog(req)
	}

	return routeNotFound(req, path, method)
}

func postDog(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var dog Dog
	err := json.Unmarshal([]byte(req.Body), &dog)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error":"%s"}`, err.Error()),
		}, nil
	}

	err = storeDog(&dog)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error":"%s"}`, err.Error()),
		}, nil
	}

	returnedDogJson, err := json.Marshal(dog)
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

func storeDog(dog *Dog) error {
	svc := dynamodb.New(session.Must(createSession()))
	dog.ID = uuid.New().String()

	tableName := os.Getenv("DOGS_TABLE_NAME")
	fmt.Printf("Storing dog in table %s\n", tableName)

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(dog.ID),
			},
			"name": {
				S: aws.String(dog.Name),
			},
		},
	}

	_, err := svc.PutItem(input)
	if err != nil {
		log.Printf("Error persisting dog to DynamoDB: %s", err)
		return err
	}

	return nil
}

func listDogs() (events.APIGatewayV2HTTPResponse, error) {
	svc := dynamodb.New(session.Must(createSession()))

	input := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DOGS_TABLE_NAME")),
		Limit:     aws.Int64(5),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf(`{"error":"Failed to scan dogs table: %s"}`, err.Error()),
		}, nil
	}

	if len(result.Items) == 0 {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       "[]",
		}, nil
	}

	var dogs []Dog
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dogs)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf(`{"error":"Failed to unmarshal dogs: %s"}`, err.Error()),
		}, nil
	}

	dogsJSON, err := json.Marshal(dogs)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf(`{"error":"Failed to marshal dogs to JSON: %s"}`, err.Error()),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(dogsJSON),
	}, nil
}

func createSession() (*session.Session, error) {
	useLocalStack := os.Getenv("USE_LOCALSTACK") == "true"
	fmt.Printf("Using localstack: %t\n", useLocalStack)

	config := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}
	fmt.Printf("Creating config for region %s\n", *config.Region)

	if useLocalStack {
		fmt.Printf("Setting endpoint to %s\n", os.Getenv("AWS_ENDPOINT_URL"))
		config.Endpoint = aws.String(os.Getenv("AWS_ENDPOINT_URL"))
		config.Credentials = credentials.NewStaticCredentials("test", "test", "")
		config.DisableSSL = aws.Bool(true)
	}

	return session.NewSession(config)
}

func healthCheck() (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Body: "OK",
	}, nil
}

func routeNotFound(req events.APIGatewayV2HTTPRequest, path string, method string) (events.APIGatewayV2HTTPResponse, error) {
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
