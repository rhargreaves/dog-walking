package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rhargreaves/dog-walking/photo-moderator/common"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.S3Event) error {
	for _, record := range req.Records {
		sourceBucket := record.S3.Bucket.Name
		sourceKey := record.S3.Object.Key
		fmt.Printf("Source bucket: %s, source key: %s\n", sourceBucket, sourceKey)

		approvedDogPhotosBucket := os.Getenv("DOG_IMAGES_BUCKET")
		dogTableName := os.Getenv("DOGS_TABLE_NAME")

		err := approvePhoto(sourceBucket, sourceKey, approvedDogPhotosBucket, dogTableName)
		if err != nil {
			fmt.Printf("Error approving photo: %s\n", err)
			return err
		}

	}
	return nil
}

func approvePhoto(pendingPhotosBucket string, dogId string, approvedPhotosBucket string,
	dogTableName string) error {
	photoHash, err := moveS3Object(pendingPhotosBucket, dogId, approvedPhotosBucket)
	if err != nil {
		fmt.Printf("Error moving object to approved bucket: %s\n", err)
		return err
	}

	err = updatePhotoStatus(dogId, photoHash, dogTableName)
	if err != nil {
		fmt.Printf("Error updating photo status: %s\n", err)
		return err
	}

	return nil
}

func updatePhotoStatus(dogId string, photoHash string, dogTableName string) error {
	dynamodbSvc := dynamodb.New(session.Must(common.CreateSession()))
	_, err := dynamodbSvc.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(dogTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(dogId),
			},
		},
		UpdateExpression: aws.String("SET photoStatus = :photoStatus, photoHash = :photoHash"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":photoStatus": {
				S: aws.String("approved"),
			},
			":photoHash": {
				S: aws.String(photoHash),
			},
		},
	})
	return err
}

func moveS3Object(sourceBucket string, sourceKey string, destinationBucket string) (string, error) {
	hash, err := copyS3Object(sourceBucket, sourceKey, destinationBucket)
	if err != nil {
		fmt.Printf("Error copying object to destination bucket: %s\n", err)
		return "", err
	}

	err = deleteS3Object(sourceBucket, sourceKey)
	if err != nil {
		fmt.Printf("Error deleting object from source bucket: %s\n", err)
		return "", err
	}
	return hash, nil
}

func copyS3Object(sourceBucket string, sourceKey string, destinationBucket string) (string, error) {
	s3svc := s3.New(session.Must(common.CreateS3Session()))
	res, err := s3svc.CopyObject(&s3.CopyObjectInput{
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

func deleteS3Object(bucket string, key string) error {
	s3svc := s3.New(session.Must(common.CreateS3Session()))
	_, err := s3svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
