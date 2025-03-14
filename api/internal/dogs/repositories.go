package dogs

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

var ErrDogNotFound = errors.New("dog not found")

type DogRepository interface {
	Create(dog *Dog) error
	List() ([]Dog, error)
	Get(id string) (*Dog, error)
	Update(id string, dog *Dog) error
}

type dogRepository struct {
	tableName string
	dynamoDB  *dynamodb.DynamoDB
}

func NewDogRepository(tableName string) DogRepository {
	dynamoDB := dynamodb.New(session.Must(createSession()))
	return &dogRepository{tableName: tableName, dynamoDB: dynamoDB}
}

func (r *dogRepository) Create(dog *Dog) error {
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

	_, err := r.dynamoDB.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put dog: %w", err)
	}
	return nil
}

func (r *dogRepository) List() ([]Dog, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := r.dynamoDB.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan dogs table: %w", err)
	}

	var dogs []Dog
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dogs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dogs: %w", err)
	}

	return dogs, nil
}

func (r *dogRepository) Get(id string) (*Dog, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	result, err := r.dynamoDB.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get dog: %w", err)
	}

	if result.Item == nil {
		return nil, ErrDogNotFound
	}

	var dog Dog
	err = dynamodbattribute.UnmarshalMap(result.Item, &dog)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dog: %w", err)
	}

	return &dog, nil
}

func (r *dogRepository) Update(id string, dog *Dog) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#n": aws.String("name"),
		},
		UpdateExpression: aws.String("set #n = :name"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(dog.Name),
			},
		},
		ReturnValues: aws.String("ALL_OLD"),
	}

	result, err := r.dynamoDB.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update dog: %w", err)
	}

	if len(result.Attributes) == 0 {
		return ErrDogNotFound
	}

	return nil
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
