package common

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func IsLocal() bool {
	return os.Getenv("USE_LOCALSTACK") == "true"
}

func CreateSession() (*session.Session, error) {
	return createSessionWithEndpoint(os.Getenv("AWS_ENDPOINT_URL"))
}

func CreateS3Session() (*session.Session, error) {
	return createSessionWithEndpoint(os.Getenv("AWS_S3_ENDPOINT_URL"))
}

func createSessionWithEndpoint(localEndpoint string) (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	config := &aws.Config{
		Region: &region,
	}
	if IsLocal() {
		config.Endpoint = aws.String(localEndpoint)
		config.Credentials = credentials.NewStaticCredentials("test", "test", "")
		config.DisableSSL = aws.Bool(true)
	}
	return session.NewSession(config)
}
