package aws_session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rhargreaves/dog-walking/shared/env"
)

func CreateS3Session() (*session.Session, error) {
	region := env.MustGetenv("AWS_REGION")
	if IsLocal() {
		return session.NewSession(&aws.Config{
			Region:      &region,
			Endpoint:    aws.String(env.MustGetenv("AWS_S3_ENDPOINT_URL")),
			Credentials: credentials.NewStaticCredentials("test", "test", ""),
		})
	}
	return session.NewSession(&aws.Config{
		Region: &region,
	})
}

func CreateSession() (*session.Session, error) {
	region := env.MustGetenv("AWS_REGION")
	if IsLocal() {
		return session.NewSession(&aws.Config{
			Region:      &region,
			Endpoint:    aws.String(env.MustGetenv("AWS_ENDPOINT_URL")),
			Credentials: credentials.NewStaticCredentials("test", "test", ""),
		})
	}
	return session.NewSession(&aws.Config{
		Region: &region,
	})
}
