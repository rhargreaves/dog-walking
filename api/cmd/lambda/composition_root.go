package main

import (
	"log"

	"github.com/rhargreaves/dog-walking/api/internal/dogs"
	"github.com/rhargreaves/dog-walking/shared/aws_session"
	"github.com/rhargreaves/dog-walking/shared/env"
)

func createHandlers() (dogs.DogHandler, dogs.DogPhotoHandler) {
	session, err := aws_session.CreateSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	dogRepository := dogs.NewDynamoDBDogRepository(dogs.DynamoDBDogRepositoryConfig{
		TableName: env.MustGetenv("DOGS_TABLE_NAME"),
	}, session)
	dogHandler := dogs.NewDogHandler(dogs.DogHandlerConfig{
		ImagesCdnBaseUrl: env.MustGetenv("CLOUDFRONT_BASE_URL"),
	}, dogRepository)

	s3session, err := aws_session.CreateS3Session()
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
