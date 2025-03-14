package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rhargreaves/dog-walking/api/models"
)

func postDog(c *gin.Context) {
	var dog models.Dog
	if err := c.ShouldBindJSON(&dog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := storeDog(&dog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dog)
}

func storeDog(dog *models.Dog) error {
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

func listDogs(c *gin.Context) {
	svc := dynamodb.New(session.Must(createSession()))

	input := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DOGS_TABLE_NAME")),
		Limit:     aws.Int64(5),
	}

	result, err := svc.Scan(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to scan dogs table: %s", err.Error())})
		return
	}

	if len(result.Items) == 0 {
		c.JSON(http.StatusOK, []models.Dog{})
		return
	}

	var dogs []models.Dog
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dogs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmarshal dogs: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, dogs)
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
