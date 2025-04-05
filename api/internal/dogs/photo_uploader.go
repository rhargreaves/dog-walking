package dogs

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type DogPhotoUploader interface {
	Upload(id string, fileData io.Reader, contentType string) (string, error)
}

type S3PhotoUploaderConfig struct {
	BucketName string
}

type s3DogPhotoUploader struct {
	config  *S3PhotoUploaderConfig
	session *session.Session
}

func NewDogPhotoUploader(s3PhotoUploaderConfig S3PhotoUploaderConfig, session *session.Session) DogPhotoUploader {
	return &s3DogPhotoUploader{config: &s3PhotoUploaderConfig, session: session}
}

func (u *s3DogPhotoUploader) Upload(id string, fileData io.Reader, contentType string) (string, error) {
	uploader := s3manager.NewUploader(u.session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(u.config.BucketName),
		Key:         aws.String(id),
		Body:        fileData,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload photo to S3: %w", err)
	}

	return strings.Replace(*result.ETag, "\"", "", 2), nil
}
