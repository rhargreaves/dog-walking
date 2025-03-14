package dogs

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type DogRepository interface {
	Create(dog *Dog) error
	List() ([]Dog, error)
}

type dogRepository struct {
	tableName string
}

func NewDogRepository(tableName string) DogRepository {
	return &dogRepository{tableName: tableName}
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

func (r *dogRepository) Create(dog *Dog) error {
	svc := dynamodb.New(session.Must(createSession()))
	dog.ID = uuid.New().String()

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
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

func (r *dogRepository) List() ([]Dog, error) {
	svc := dynamodb.New(session.Must(createSession()))

	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	var dogs []Dog
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dogs)
	if err != nil {
		return nil, err
	}
	return dogs, nil
}
