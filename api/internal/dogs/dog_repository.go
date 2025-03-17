package dogs

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/models"
)

var ErrDogNotFound = errors.New("dog not found")

//go:generate mockery --name DogRepository --output ../mocks --outpkg mocks --case underscore
type DogRepository interface {
	Create(dog *models.Dog) error
	List() ([]models.Dog, error)
	Get(id string) (*models.Dog, error)
	Update(id string, dog *models.Dog) error
	Delete(id string) error
}

type dogRepository struct {
	tableName string
	dynamoDB  *dynamodb.DynamoDB
}

func NewDogRepository(tableName string) DogRepository {
	dynamoDB := dynamodb.New(session.Must(common.CreateSession()))
	return &dogRepository{tableName: tableName, dynamoDB: dynamoDB}
}

func (r *dogRepository) Create(dog *models.Dog) error {
	dog.ID = uuid.New().String()

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id":   {S: aws.String(dog.ID)},
			"name": {S: aws.String(dog.Name)},
		},
	}

	_, err := r.dynamoDB.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put dog: %w", err)
	}
	return nil
}

func (r *dogRepository) List() ([]models.Dog, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := r.dynamoDB.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan dogs table: %w", err)
	}

	var dogs []models.Dog
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dogs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dogs: %w", err)
	}

	return dogs, nil
}

func (r *dogRepository) Get(id string) (*models.Dog, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       createKey(id),
	}

	result, err := r.dynamoDB.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get dog: %w", err)
	}

	if result.Item == nil {
		return nil, ErrDogNotFound
	}

	var dog models.Dog
	err = dynamodbattribute.UnmarshalMap(result.Item, &dog)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dog: %w", err)
	}

	return &dog, nil
}

func (r *dogRepository) Update(id string, dog *models.Dog) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key:       createKey(id),
		ExpressionAttributeNames: map[string]*string{
			"#n": aws.String("name"),
		},
		UpdateExpression: aws.String("set #n = :name, breed = :breed"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(dog.Name),
			},
			":breed": {
				S: aws.String(dog.Breed),
			},
		},
		ConditionExpression: aws.String("attribute_exists(id)"),
	}

	_, err := r.dynamoDB.UpdateItem(input)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() ==
			dynamodb.ErrCodeConditionalCheckFailedException {
			return ErrDogNotFound
		}
		return fmt.Errorf("failed to update dog: %w", err)
	}

	return nil
}

func (r *dogRepository) Delete(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName:    aws.String(r.tableName),
		Key:          createKey(id),
		ReturnValues: aws.String("ALL_OLD"),
	}

	result, err := r.dynamoDB.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("failed to delete dog: %w", err)
	}

	if result.Attributes == nil {
		return ErrDogNotFound
	}

	return nil
}

func createKey(id string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}
}
