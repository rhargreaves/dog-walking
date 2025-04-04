package mocks

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	UpdateItemFunc func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
}

func (t *MockDynamoDB) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return t.UpdateItemFunc(input)
}

type MockS3 struct {
	s3iface.S3API
	PutObjectFunc    func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObjectFunc func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	CopyObjectFunc   func(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error)
}

func (t *MockS3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return t.PutObjectFunc(input)
}

func (t *MockS3) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return t.DeleteObjectFunc(input)
}

func (t *MockS3) CopyObject(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	return t.CopyObjectFunc(input)
}
