package dogs

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rhargreaves/dog-walking/api/internal/common"
)

type DogPhotoRepository interface {
	Upload(id string, fileData io.Reader, contentType string) error
}

type dogPhotoRepository struct {
	bucketName string
}

func NewDogPhotoRepository(bucketName string) DogPhotoRepository {
	return &dogPhotoRepository{bucketName: bucketName}
}

func (r *dogPhotoRepository) Upload(id string, fileData io.Reader, contentType string) error {
	sess, err := common.CreateS3Session()
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(id),
		Body:        fileData,
		ContentType: aws.String(contentType),
	})

	return err
}
