package dogs

import (
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rhargreaves/dog-walking/api/internal/common"
)

type DogPhotoUploader interface {
	Upload(id string, fileData io.Reader, contentType string) error
}

type S3PhotoUploaderConfig struct {
	BucketName    string
	IsLocal       bool
	LocalEndpoint string
	Region        string
}

type s3DogPhotoUploader struct {
	config        *S3PhotoUploaderConfig
	dogRepository DogRepository
}

func NewDogPhotoUploader(s3PhotoUploaderConfig S3PhotoUploaderConfig, dogRepository DogRepository) DogPhotoUploader {
	return &s3DogPhotoUploader{config: &s3PhotoUploaderConfig, dogRepository: dogRepository}
}

func (r *s3DogPhotoUploader) Upload(id string, fileData io.Reader, contentType string) error {
	sess, err := common.CreateSession(r.config.IsLocal, r.config.LocalEndpoint, r.config.Region)
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(r.config.BucketName),
		Key:         aws.String(id),
		Body:        fileData,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return err
	}

	photoHash := strings.Replace(*result.ETag, "\"", "", 2)
	err = r.dogRepository.UpdatePhotoHash(id, photoHash)
	if err != nil {
		return err
	}
	return err
}
