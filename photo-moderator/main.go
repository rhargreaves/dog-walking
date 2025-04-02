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

func approvePhoto(dogId string, photoHash string) error {
	dogTable := os.Getenv("DOGS_TABLE_NAME")
	dynamodbSvc := dynamodb.New(session.Must(common.CreateSession()))
	_, err := dynamodbSvc.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(dogTable),
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

func handler(ctx context.Context, req events.S3Event) error {
	for _, record := range req.Records {
		sourceBucket := record.S3.Bucket.Name
		sourceKey := record.S3.Object.Key
		fmt.Printf("Source bucket: %s, source key: %s\n", sourceBucket, sourceKey)

		// move image to the dog-photos bucket
		approvedDogPhotosBucket := os.Getenv("DOG_IMAGES_BUCKET")
		fmt.Printf("Approved dog photos bucket: %s\n", approvedDogPhotosBucket)

		// move the object to the dog-photos bucket
		s3svc := s3.New(session.Must(common.CreateS3Session()))
		res, err := s3svc.CopyObject(&s3.CopyObjectInput{
			Bucket:     aws.String(approvedDogPhotosBucket),
			CopySource: aws.String(fmt.Sprintf("%s/%s", sourceBucket, sourceKey)),
			Key:        aws.String(sourceKey),
		})
		if err != nil {
			fmt.Printf("Error moving object to approved bucket: %s\n", err)
			return err
		}
		photoHash := strings.Trim(*res.CopyObjectResult.ETag, "\"")

		// delete the object from the pending-dog-images bucket
		_, err = s3svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(sourceBucket),
			Key:    aws.String(sourceKey),
		})
		if err != nil {
			fmt.Printf("Error deleting object from pending bucket: %s\n", err)
			return err
		}

		err = approvePhoto(sourceKey, photoHash)
		if err != nil {
			fmt.Printf("Error updating photo status: %s\n", err)
			return err
		}
	}
	return nil
}
