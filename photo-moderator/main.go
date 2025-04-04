package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	"github.com/rhargreaves/dog-walking/photo-moderator/common"
	"github.com/rhargreaves/dog-walking/photo-moderator/moderator"
	"github.com/rhargreaves/dog-walking/photo-moderator/moderator/rekognition_stub"
)

func main() {
	lambda.Start(handler)
}

func newRekognitionClient(isLocal bool, session *session.Session, s3session *session.Session) rekognitioniface.RekognitionAPI {
	if isLocal {
		return rekognition_stub.NewStubRekognitionClient(s3session)
	}
	return rekognition.New(session)
}

func handler(ctx context.Context, req events.S3Event) error {
	for _, record := range req.Records {
		sourceBucket := record.S3.Bucket.Name
		sourceKey := record.S3.Object.Key
		fmt.Printf("Source bucket: %s, source key: %s\n", sourceBucket, sourceKey)

		approvedDogPhotosBucket := os.Getenv("DOG_IMAGES_BUCKET")
		dogTableName := os.Getenv("DOGS_TABLE_NAME")

		dbSession := session.Must(common.CreateSession())
		s3session := session.Must(common.CreateS3Session())

		rekClient := newRekognitionClient(common.IsLocal(), dbSession, s3session)
		breedDetector := moderator.NewBreedDetector(moderator.BreedDetectorConfig{
			BucketName: sourceBucket,
		}, rekClient)

		moderator := moderator.NewModerator(dogTableName, approvedDogPhotosBucket, breedDetector, dbSession, s3session)
		err := moderator.ModeratePhoto(sourceBucket, sourceKey)
		if err != nil {
			fmt.Printf("Error approving photo: %s\n", err)
			return err
		}

	}
	return nil
}
