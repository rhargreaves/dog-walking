package moderator

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	aws_clients "github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/aws"
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
	dynamodbClient       aws_clients.DynamoDBClient
	s3Client             aws_clients.S3Client
	contentScreener      content_screener.ContentScreener
}

func NewModerator(dogTableName string, approvedPhotosBucket string, breedDetector breed_detector.BreedDetector,
	dynamodbClient aws_clients.DynamoDBClient, s3Client aws_clients.S3Client,
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
		fmt.Printf("Error screening image: %s\n", err)
		return "", err
	}
	if !contentScreenerResult.IsSafe {
		err = m.rejectDogPhoto(dogId)
		if err != nil {
			fmt.Printf("Error rejecting dog photo: %s\n", err)
			return "", err
		}
		fmt.Printf("Dog photo rejected: %s\n", contentScreenerResult.Reason)
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
			fmt.Printf("Error rejecting dog photo: %s\n", err)
			return "", err
		}
		return PhotoStatusRejected, nil
	case breed_detector.ErrNoSpecificBreedDetected:
		err = m.approveDogPhoto(dogId, nil, pendingPhotosBucket)
		if err != nil {
			fmt.Printf("Error approving dog photo: %s\n", err)
			return "", err
		}
		return PhotoStatusApproved, nil
	default:
		fmt.Printf("Error moderating photo: %s\n", err)
		return "", err
	}
}

func (m *moderator) approveDogPhoto(dogId string, breed *string, pendingPhotosBucket string) error {
	photoHash, err := m.moveS3Object(pendingPhotosBucket, dogId, m.approvedPhotosBucket)
	if err != nil {
		fmt.Printf("Error moving object to approved bucket: %s\n", err)
		return err
	}

	if breed != nil {
		err = m.updateDogRecord(dogId, photoHash, *breed)
	} else {
		err = m.updateDogRecordWithoutBreed(dogId, photoHash)
	}
	if err != nil {
		fmt.Printf("Error updating dog record: %s\n", err)
		return err
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
	return err
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
	return err
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
