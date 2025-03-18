package common

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func IsLocal() bool {
	return os.Getenv("USE_LOCALSTACK") == "true"
}

func CreateSession() (*session.Session, error) {
	fmt.Printf("Using localstack: %t\n", IsLocal())

	config := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}
	fmt.Printf("Creating config for region %s\n", *config.Region)

	if IsLocal() {
		config.Endpoint = aws.String(os.Getenv("AWS_ENDPOINT_URL"))
		config.Credentials = credentials.NewStaticCredentials("test", "test", "")
		config.DisableSSL = aws.Bool(true)
		fmt.Printf("Setting endpoint to %s\n", *config.Endpoint)
	}

	return session.NewSession(config)
}

func CreateS3Session() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	if IsLocal() {
		return session.NewSession(&aws.Config{
			Region:      &region,
			Endpoint:    aws.String(os.Getenv("AWS_S3_ENDPOINT_URL")),
			Credentials: credentials.NewStaticCredentials("test", "test", ""),
		})
	}
	return session.NewSession(&aws.Config{
		Region: &region,
	})
}
