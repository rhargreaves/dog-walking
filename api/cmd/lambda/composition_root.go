package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rhargreaves/dog-walking/api/internal/dogs"
	"github.com/rhargreaves/dog-walking/shared/env"
)

func createHandlers(isLocal bool) (dogs.DogHandler, dogs.DogPhotoHandler) {
	session, err := createSession(isLocal, false)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	dogRepository := dogs.NewDynamoDBDogRepository(dogs.DynamoDBDogRepositoryConfig{
		TableName: env.MustGetenv("DOGS_TABLE_NAME"),
	}, session)
	dogHandler := dogs.NewDogHandler(dogs.DogHandlerConfig{
		ImagesCdnBaseUrl: env.MustGetenv("CLOUDFRONT_BASE_URL"),
	}, dogRepository)

	s3session, err := createSession(isLocal, true)
	if err != nil {
		log.Fatalf("Failed to create S3 session: %v", err)
	}
	dogImagesBucket := env.MustGetenv("PENDING_DOG_IMAGES_BUCKET")
	dogPhotoUploader := dogs.NewDogPhotoUploader(dogs.S3PhotoUploaderConfig{
		BucketName: dogImagesBucket,
	}, s3session)
	dogPhotoHandler := dogs.NewDogPhotoHandler(dogRepository, dogPhotoUploader)

	return dogHandler, dogPhotoHandler
}

func createSession(isLocal bool, forS3 bool) (*session.Session, error) {
	region := env.MustGetenv("AWS_REGION")
	config := &aws.Config{
		Region: &region,
	}
	if !isLocal {
		return session.NewSession(config)
	}

	var endpoint string
	if forS3 {
		endpoint = env.MustGetenv("AWS_S3_ENDPOINT_URL")
	} else {
		endpoint = env.MustGetenv("AWS_ENDPOINT_URL")
	}
	config.Endpoint = aws.String(endpoint)
	config.Credentials = credentials.NewStaticCredentials("test", "test", "")
	config.DisableSSL = aws.Bool(true)
	return session.NewSession(config)
}
