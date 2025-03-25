package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateSession(isLocal bool, localEndpoint string, region string) (*session.Session, error) {
	config := &aws.Config{
		Region: &region,
	}
	if isLocal {
		config.Endpoint = aws.String(localEndpoint)
		config.Credentials = credentials.NewStaticCredentials("test", "test", "")
		config.DisableSSL = aws.Bool(true)
	}
	return session.NewSession(config)
}
