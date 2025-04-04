package moderator

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	aws_clients "github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/aws"
	breed_detector "github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/breed_detector"
)

const (
	PhotoStatusApproved = "approved"
	PhotoStatusRejected = "rejected"
)

type Moderator interface {
	ModeratePhoto(pendingPhotosBucket string, dogId string) error
}

type moderator struct {
	dogTableName         string
	approvedPhotosBucket string
	breedDetector        breed_detector.BreedDetector
	dynamodbClient       aws_clients.DynamoDBClient
	s3Client             aws_clients.S3Client
}

func NewModerator(dogTableName string, approvedPhotosBucket string, breedDetector breed_detector.BreedDetector,
	dynamodbClient aws_clients.DynamoDBClient, s3Client aws_clients.S3Client) Moderator {
	return &moderator{
		dogTableName:         dogTableName,
		approvedPhotosBucket: approvedPhotosBucket,
		breedDetector:        breedDetector,
		dynamodbClient:       dynamodbClient,
		s3Client:             s3Client,
	}
}

func (m *moderator) ModeratePhoto(pendingPhotosBucket string, dogId string) error {
	breedDetectionResult, err := m.breedDetector.DetectBreed(dogId)
	if err != nil {

		if err == breed_detector.ErrNoDogDetected || err == breed_detector.ErrNoSpecificBreedDetected {
			err = m.updateDogRecordToRejected(dogId)
			if err != nil {
				fmt.Printf("Error updating dog record: %s\n", err)
				return err
			}
			return nil
		}

		fmt.Printf("Error detecting breed: %s\n", err)
		return err
	}

	if breedDetectionResult.Breed != "" {
		photoHash, err := m.moveS3Object(pendingPhotosBucket, dogId, m.approvedPhotosBucket)
		if err != nil {
			fmt.Printf("Error moving object to approved bucket: %s\n", err)
			return err
		}

		err = m.updateDogRecord(dogId, photoHash, m.dogTableName, breedDetectionResult.Breed)
		if err != nil {
			fmt.Printf("Error updating dog record: %s\n", err)
			return err
		}
	}

	return nil
}

func (m *moderator) updateDogRecordToRejected(dogId string) error {
	_, err := m.dynamodbClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(m.dogTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(dogId),
			},
		},
		UpdateExpression: aws.String("SET photoStatus = :photoStatus"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":photoStatus": {
				S: aws.String(PhotoStatusRejected),
			},
		},
	})
	return err
}

func (m *moderator) updateDogRecord(dogId string, photoHash string, dogTableName string, breed string) error {
	_, err := m.dynamodbClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(dogTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(dogId),
			},
		},
		UpdateExpression: aws.String("SET photoStatus = :photoStatus, photoHash = :photoHash, breed = :breed"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":photoStatus": {
				S: aws.String(PhotoStatusApproved),
			},
			":photoHash": {
				S: aws.String(photoHash),
			},
			":breed": {
				S: aws.String(breed),
			},
		},
	})
	return err
}

func (m *moderator) moveS3Object(sourceBucket string, sourceKey string, destinationBucket string) (string, error) {
	hash, err := m.copyS3Object(sourceBucket, sourceKey, destinationBucket)
	if err != nil {
		fmt.Printf("Error copying object to destination bucket: %s\n", err)
		return "", err
	}

	err = m.deleteS3Object(sourceBucket, sourceKey)
	if err != nil {
		fmt.Printf("Error deleting object from source bucket: %s\n", err)
		return "", err
	}
	return hash, nil
}

func (m *moderator) copyS3Object(sourceBucket string, sourceKey string, destinationBucket string) (string, error) {
	res, err := m.s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(destinationBucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", sourceBucket, sourceKey)),
		Key:        aws.String(sourceKey),
	})
	if err != nil {
		fmt.Printf("Error copying object to destination bucket: %s\n", err)
		return "", err
	}
	hash := strings.Trim(*res.CopyObjectResult.ETag, "\"")
	return hash, nil
}

func (m *moderator) deleteS3Object(bucket string, key string) error {
	_, err := m.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
