package dogs

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rhargreaves/dog-walking/api/internal/common"
)

func PostDog(c *gin.Context) {
	var dog Dog
	if err := c.ShouldBindJSON(&dog); err != nil {
		c.Error(common.APIError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err := storeDog(&dog)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dog)
}

func storeDog(dog *Dog) error {
	svc := dynamodb.New(session.Must(createSession()))
	dog.ID = uuid.New().String()

	tableName := os.Getenv("DOGS_TABLE_NAME")

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
		return err
	}
	return nil
}

func ListDogs(c *gin.Context) {
	svc := dynamodb.New(session.Must(createSession()))

	input := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DOGS_TABLE_NAME")),
	}

	result, err := svc.Scan(input)
	if err != nil {
		c.Error(fmt.Errorf("failed to scan dogs table: %s", err.Error()))
		return
	}

	var dogs []Dog
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dogs)
	if err != nil {
		c.Error(fmt.Errorf("failed to unmarshal dogs: %s", err.Error()))
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
		config.Endpoint = aws.String(os.Getenv("AWS_ENDPOINT_URL"))
		config.Credentials = credentials.NewStaticCredentials("test", "test", "")
		config.DisableSSL = aws.Bool(true)
		fmt.Printf("Setting endpoint to %s\n", *config.Endpoint)
	}

	return session.NewSession(config)
}
