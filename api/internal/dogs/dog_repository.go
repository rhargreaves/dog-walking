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
	"github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
)

var ErrDogNotFound = errors.New("dog not found")

type DogRepository interface {
	Create(dog *domain.Dog) error
	List(limit int, nextToken string) (*domain.DogList, error)
	Get(id string) (*domain.Dog, error)
	Update(id string, dog *domain.Dog) error
	Delete(id string) error
	UpdatePhotoHash(id string, photoHash string) error
}

type DynamoDBDogRepositoryConfig struct {
	TableName string
}

type dynamoDBDogRepository struct {
	config   *DynamoDBDogRepositoryConfig
	dynamoDB *dynamodb.DynamoDB
}

func NewDynamoDBDogRepository(dynamoDBDogRepositoryConfig DynamoDBDogRepositoryConfig) DogRepository {
	dynamoDB := dynamodb.New(session.Must(common.CreateSession()))
	return &dynamoDBDogRepository{config: &dynamoDBDogRepositoryConfig, dynamoDB: dynamoDB}
}

func (r *dynamoDBDogRepository) Create(dog *domain.Dog) error {
	dog.ID = uuid.New().String()

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.config.TableName),
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

func (r *dynamoDBDogRepository) List(limit int, nextToken string) (*domain.DogList, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.config.TableName),
		Limit:     aws.Int64(int64(limit)),
	}

	if nextToken != "" {
		input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(nextToken)},
		}
	}

	result, err := r.dynamoDB.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan dogs table: %w", err)
	}

	var dogs []domain.Dog
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dogs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dogs: %w", err)
	}

	nextToken = ""
	lastEvaluatedKey := result.LastEvaluatedKey["id"]
	if lastEvaluatedKey != nil {
		nextToken = *lastEvaluatedKey.S
	}

	return &domain.DogList{
		Dogs:      dogs,
		NextToken: nextToken,
	}, nil
}

func (r *dynamoDBDogRepository) Get(id string) (*domain.Dog, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.config.TableName),
		Key:       createKey(id),
	}

	result, err := r.dynamoDB.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get dog: %w", err)
	}

	if result.Item == nil {
		return nil, ErrDogNotFound
	}

	var dog domain.Dog
	err = dynamodbattribute.UnmarshalMap(result.Item, &dog)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dog: %w", err)
	}

	return &dog, nil
}

func (r *dynamoDBDogRepository) Update(id string, dog *domain.Dog) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.config.TableName),
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

func (r *dynamoDBDogRepository) UpdatePhotoHash(id string, photoHash string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.config.TableName),
		Key:       createKey(id),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":photoHash": {S: aws.String(photoHash)},
		},
		UpdateExpression:    aws.String("set photoHash = :photoHash"),
		ConditionExpression: aws.String("attribute_exists(id)"),
	}

	_, err := r.dynamoDB.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update dog photo hash: %w", err)
	}

	return nil
}

func (r *dynamoDBDogRepository) Delete(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName:    aws.String(r.config.TableName),
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
