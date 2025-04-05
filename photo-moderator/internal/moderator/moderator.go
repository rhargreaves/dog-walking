package moderator

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	breed_detector "github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/breed_detector"
	"github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/content_screener"
)

const (
	PhotoStatusApproved = "approved"
	PhotoStatusRejected = "rejected"
)

type Moderator interface {
	ModeratePhoto(pendingPhotosBucket string, dogId string) (string, error)
}

type moderator struct {
	dogTableName         string
	approvedPhotosBucket string
	breedDetector        breed_detector.BreedDetector
	dynamodbClient       dynamodbiface.DynamoDBAPI
	s3Client             s3iface.S3API
	contentScreener      content_screener.ContentScreener
}

func NewModerator(dogTableName string, approvedPhotosBucket string, breedDetector breed_detector.BreedDetector,
	dynamodbClient dynamodbiface.DynamoDBAPI, s3Client s3iface.S3API,
	contentScreener content_screener.ContentScreener) Moderator {
	return &moderator{
		dogTableName:         dogTableName,
		approvedPhotosBucket: approvedPhotosBucket,
		breedDetector:        breedDetector,
		dynamodbClient:       dynamodbClient,
		s3Client:             s3Client,
		contentScreener:      contentScreener,
	}
}

func (m *moderator) ModeratePhoto(pendingPhotosBucket string, dogId string) (string, error) {
	contentScreenerResult, err := m.contentScreener.ScreenImage(dogId)
	if err != nil {
		return "", fmt.Errorf("failed to screen image for inappropriate content: %w", err)
	}
	if !contentScreenerResult.IsSafe {
		err = m.rejectDogPhoto(dogId)
		if err != nil {
			return "", fmt.Errorf("failed to reject dog photo after content screening: %w", err)
		}
		return PhotoStatusRejected, nil
	}

	breedDetectionResult, err := m.breedDetector.DetectBreed(dogId)
	if err == nil {
		return PhotoStatusApproved, m.approveDogPhoto(dogId, &breedDetectionResult.Breed, pendingPhotosBucket)
	}

	switch err {
	case breed_detector.ErrNoDogDetected:
		err = m.rejectDogPhoto(dogId)
		if err != nil {
			return "", fmt.Errorf("failed to reject dog photo after no dog detected: %w", err)
		}
		return PhotoStatusRejected, nil
	case breed_detector.ErrNoSpecificBreedDetected:
		err = m.approveDogPhoto(dogId, nil, pendingPhotosBucket)
		if err != nil {
			return "", fmt.Errorf("failed to approve dog photo with no specific breed: %w", err)
		}
		return PhotoStatusApproved, nil
	default:
		return "", fmt.Errorf("failed to moderate photo: %w", err)
	}
}

func (m *moderator) approveDogPhoto(dogId string, breed *string, pendingPhotosBucket string) error {
	photoHash, err := m.moveS3Object(pendingPhotosBucket, dogId, m.approvedPhotosBucket)
	if err != nil {
		return fmt.Errorf("failed to move photo to approved bucket: %w", err)
	}

	if breed != nil {
		err = m.updateDogRecord(dogId, photoHash, *breed)
	} else {
		err = m.updateDogRecordWithoutBreed(dogId, photoHash)
	}
	if err != nil {
		return fmt.Errorf("failed to update dog record after approval: %w", err)
	}
	return nil
}

func (m *moderator) rejectDogPhoto(dogId string) error {
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
	if err != nil {
		return fmt.Errorf("failed to update dog record with rejected status: %w", err)
	}
	return nil
}

func (m *moderator) updateDogRecord(dogId string, photoHash string, breed string) error {
	_, err := m.dynamodbClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(m.dogTableName),
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
	if err != nil {
		return fmt.Errorf("failed to update dog record: %w", err)
	}
	return nil
}

func (m *moderator) updateDogRecordWithoutBreed(dogId string, photoHash string) error {
	_, err := m.dynamodbClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(m.dogTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(dogId),
			},
		},
		UpdateExpression: aws.String("SET photoStatus = :photoStatus, photoHash = :photoHash"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":photoStatus": {
				S: aws.String(PhotoStatusApproved),
			},
			":photoHash": {
				S: aws.String(photoHash),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update dog record: %w", err)
	}
	return nil
}

func (m *moderator) moveS3Object(sourceBucket string, sourceKey string, destinationBucket string) (string, error) {
	hash, err := m.copyS3Object(sourceBucket, sourceKey, destinationBucket)
	if err != nil {
		return "", fmt.Errorf("failed to copy photo to destination bucket: %w", err)
	}

	err = m.deleteS3Object(sourceBucket, sourceKey)
	if err != nil {
		return "", fmt.Errorf("failed to delete photo from source bucket: %w", err)
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
		return "", fmt.Errorf("failed to copy photo to approved bucket: %w", err)
	}
	hash := strings.Trim(*res.CopyObjectResult.ETag, "\"")
	return hash, nil
}

func (m *moderator) deleteS3Object(bucket string, key string) error {
	_, err := m.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete photo from bucket: %w", err)
	}
	return nil
}
