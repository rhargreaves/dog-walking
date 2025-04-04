package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rhargreaves/dog-walking/api/internal/dogs"
)

func createHandlers(isLocal bool) (dogs.DogHandler, dogs.DogPhotoHandler) {
	session, err := createSession(isLocal, false)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	dogRepository := dogs.NewDynamoDBDogRepository(dogs.DynamoDBDogRepositoryConfig{
		TableName: mustGetenv("DOGS_TABLE_NAME"),
	}, session)
	dogHandler := dogs.NewDogHandler(dogs.DogHandlerConfig{
		ImagesCdnBaseUrl: mustGetenv("CLOUDFRONT_BASE_URL"),
	}, dogRepository)

	s3session, err := createSession(isLocal, true)
	if err != nil {
		log.Fatalf("Failed to create S3 session: %v", err)
	}
	dogImagesBucket := mustGetenv("PENDING_DOG_IMAGES_BUCKET")
	dogPhotoUploader := dogs.NewDogPhotoUploader(dogs.S3PhotoUploaderConfig{
		BucketName: dogImagesBucket,
	}, s3session)
	dogPhotoHandler := dogs.NewDogPhotoHandler(dogRepository, dogPhotoUploader)

	return dogHandler, dogPhotoHandler
}

func createSession(isLocal bool, forS3 bool) (*session.Session, error) {
	region := mustGetenv("AWS_REGION")
	config := &aws.Config{
		Region: &region,
	}
	if !isLocal {
		return session.NewSession(config)
	}

	var endpoint string
	if forS3 {
		endpoint = mustGetenv("AWS_S3_ENDPOINT_URL")
	} else {
		endpoint = mustGetenv("AWS_ENDPOINT_URL")
	}
	config.Endpoint = aws.String(endpoint)
	config.Credentials = credentials.NewStaticCredentials("test", "test", "")
	config.DisableSSL = aws.Bool(true)
	return session.NewSession(config)
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Required environment variable not set: " + key)
	}
	return val
}
