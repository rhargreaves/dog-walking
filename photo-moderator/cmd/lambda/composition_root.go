package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator"
	"github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/breed_detector"
	"github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/rekognition_stub"
	"github.com/rhargreaves/dog-walking/shared/aws_session"
	"github.com/rhargreaves/dog-walking/shared/env"
)

func rekognitionClient(isLocal bool, s3Svc *s3.S3, session *session.Session) rekognitioniface.RekognitionAPI {
	if isLocal {
		return rekognition_stub.NewStubRekognitionClient(s3Svc)
	}
	return rekognition.New(session)
}

func createModerator(sourceBucket string) moderator.Moderator {
	approvedDogPhotosBucket := env.MustGetenv("DOG_IMAGES_BUCKET")
	dogTableName := env.MustGetenv("DOGS_TABLE_NAME")

	s3session := session.Must(aws_session.CreateS3Session())
	dbSvc := dynamodb.New(session.Must(aws_session.CreateSession()))
	s3Svc := s3.New(s3session)
	rekClient := rekognitionClient(aws_session.IsLocal(), s3Svc, s3session)

	breedDetector := breed_detector.NewBreedDetector(breed_detector.BreedDetectorConfig{
		BucketName: sourceBucket,
	}, rekClient)

	return moderator.NewModerator(dogTableName, approvedDogPhotosBucket, breedDetector, dbSvc, s3Svc)
}
