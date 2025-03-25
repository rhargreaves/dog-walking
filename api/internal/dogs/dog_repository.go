package dogs

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
)

var ErrDogNotFound = errors.New("dog not found")

type DogRepository interface {
	Create(dog *domain.Dog) error
	List(limit int, name string, nextToken string) (*domain.DogList, error)
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

func NewDynamoDBDogRepository(dynamoDBDogRepositoryConfig DynamoDBDogRepositoryConfig, session *session.Session) DogRepository {
	dynamoDB := dynamodb.New(session)
	return &dynamoDBDogRepository{config: &dynamoDBDogRepositoryConfig, dynamoDB: dynamoDB}
}

func (r *dynamoDBDogRepository) Create(dog *domain.Dog) error {
	dog.ID = uuid.New().String()

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.config.TableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id":          {S: aws.String(dog.ID)},
			"name":        {S: aws.String(dog.Name)},
			"breed":       {S: aws.String(dog.Breed)},
			"sex":         {S: aws.String(dog.Sex)},
			"isNeutered":  {BOOL: aws.Bool(dog.IsNeutered)},
			"energyLevel": {N: aws.String(strconv.Itoa(dog.EnergyLevel))},
			"size":        {S: aws.String(dog.Size)},
			"socialization": {M: map[string]*dynamodb.AttributeValue{
				"goodWithChildren":  {BOOL: aws.Bool(dog.Socialization.GoodWithChildren)},
				"goodWithPuppies":   {BOOL: aws.Bool(dog.Socialization.GoodWithPuppies)},
				"goodWithLargeDogs": {BOOL: aws.Bool(dog.Socialization.GoodWithLargeDogs)},
				"goodWithSmallDogs": {BOOL: aws.Bool(dog.Socialization.GoodWithSmallDogs)},
			}},
			"specialInstructions": {S: aws.String(dog.SpecialInstructions)},
			"dateOfBirth":         {S: aws.String(dog.DateOfBirth)},
		},
	}

	_, err := r.dynamoDB.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put dog: %w", err)
	}
	return nil
}

func (r *dynamoDBDogRepository) List(limit int, name string, nextToken string) (*domain.DogList, error) {
	var dogs []domain.Dog
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue
	var lastProcessedKey map[string]*dynamodb.AttributeValue

	if nextToken != "" {
		lastEvaluatedKey = map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(nextToken)},
		}
	}

	// Keep scanning until we have enough items or reach the end of the table
	for len(dogs) < limit {
		input := &dynamodb.ScanInput{
			TableName:         aws.String(r.config.TableName),
			ExclusiveStartKey: lastEvaluatedKey,
			Limit:             aws.Int64(25),
		}

		if name != "" {
			input.ExpressionAttributeNames = map[string]*string{
				"#n": aws.String("name"),
			}
			input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
				":name": {S: aws.String(name)},
			}
			input.FilterExpression = aws.String("contains(#n, :name)")
		}

		result, err := r.dynamoDB.Scan(input)
		if err != nil {
			return nil, fmt.Errorf("failed to scan dogs table: %w", err)
		}

		var batchDogs []domain.Dog
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &batchDogs)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal dogs: %w", err)
		}

		// Process items one by one to maintain accurate pagination
		for _, dog := range batchDogs {
			if len(dogs) >= limit {
				// Store the last processed key for the next page
				lastProcessedKey = map[string]*dynamodb.AttributeValue{
					"id": {S: aws.String(dog.ID)},
				}
				break
			}
			dogs = append(dogs, dog)
		}

		lastEvaluatedKey = result.LastEvaluatedKey
		if lastEvaluatedKey == nil {
			break
		}
	}

	// Set the next token based on the last processed item
	nextToken = ""
	if lastProcessedKey != nil {
		nextToken = *lastProcessedKey["id"].S
	} else if lastEvaluatedKey != nil {
		nextToken = *lastEvaluatedKey["id"].S
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
		UpdateExpression: aws.String("set #n = :name," +
			"breed = :breed," +
			"sex = :sex," +
			"isNeutered = :isNeutered," +
			"energyLevel = :energyLevel," +
			"size = :size," +
			"socialization = :socialization," +
			"specialInstructions = :specialInstructions," +
			"dateOfBirth = :dateOfBirth"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(dog.Name),
			},
			":breed": {
				S: aws.String(dog.Breed),
			},
			":sex": {
				S: aws.String(dog.Sex),
			},
			":isNeutered": {
				BOOL: aws.Bool(dog.IsNeutered),
			},
			":energyLevel": {
				N: aws.String(strconv.Itoa(dog.EnergyLevel)),
			},
			":size": {
				S: aws.String(dog.Size),
			},
			":socialization": {
				M: map[string]*dynamodb.AttributeValue{
					"goodWithChildren":  {BOOL: aws.Bool(dog.Socialization.GoodWithChildren)},
					"goodWithPuppies":   {BOOL: aws.Bool(dog.Socialization.GoodWithPuppies)},
					"goodWithLargeDogs": {BOOL: aws.Bool(dog.Socialization.GoodWithLargeDogs)},
					"goodWithSmallDogs": {BOOL: aws.Bool(dog.Socialization.GoodWithSmallDogs)},
				},
			},
			":specialInstructions": {
				S: aws.String(dog.SpecialInstructions),
			},
			":dateOfBirth": {
				S: aws.String(dog.DateOfBirth),
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
