package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rhargreaves/dog-walking/photo-moderator/common"
	"github.com/rhargreaves/dog-walking/photo-moderator/moderator"
	"github.com/rhargreaves/dog-walking/photo-moderator/moderator/rekognition_stub"
)

func rekognitionClient(isLocal bool, s3Svc *s3.S3, session *session.Session) rekognitioniface.RekognitionAPI {
	if isLocal {
		return rekognition_stub.NewStubRekognitionClient(s3Svc)
	}
	return rekognition.New(session)
}

func createModerator(sourceBucket string) moderator.Moderator {
	approvedDogPhotosBucket := os.Getenv("DOG_IMAGES_BUCKET")
	dogTableName := os.Getenv("DOGS_TABLE_NAME")

	s3session := session.Must(common.CreateS3Session())
	dbSvc := dynamodb.New(session.Must(common.CreateSession()))
	s3Svc := s3.New(s3session)
	rekClient := rekognitionClient(common.IsLocal(), s3Svc, s3session)

	breedDetector := moderator.NewBreedDetector(moderator.BreedDetectorConfig{
		BucketName: sourceBucket,
	}, rekClient)

	return moderator.NewModerator(dogTableName, approvedDogPhotosBucket, breedDetector, dbSvc, s3Svc)
}
