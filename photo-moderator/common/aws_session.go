package common

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

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

func CreateSession() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	if IsLocal() {
		return session.NewSession(&aws.Config{
			Region:      &region,
			Endpoint:    aws.String(os.Getenv("AWS_ENDPOINT_URL")),
			Credentials: credentials.NewStaticCredentials("test", "test", ""),
		})
	}
	return session.NewSession(&aws.Config{
		Region: &region,
	})
}
