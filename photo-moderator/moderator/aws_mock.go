package moderator

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

type TestDynamoDB struct {
	DynamoDBClient
	UpdateItemFunc func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
}

func (t *TestDynamoDB) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return t.UpdateItemFunc(input)
}

type TestS3 struct {
	S3Client
	PutObjectFunc    func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObjectFunc func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	CopyObjectFunc   func(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error)
}

func (t *TestS3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return t.PutObjectFunc(input)
}

func (t *TestS3) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return t.DeleteObjectFunc(input)
}

func (t *TestS3) CopyObject(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	return t.CopyObjectFunc(input)
}
