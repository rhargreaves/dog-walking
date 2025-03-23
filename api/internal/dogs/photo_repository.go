package dogs

import (
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rhargreaves/dog-walking/api/internal/common"
)

type DogPhotoRepository interface {
	Upload(id string, fileData io.Reader, contentType string) error
}

type dogPhotoRepository struct {
	bucketName    string
	dogRepository DogRepository
}

func NewDogPhotoRepository(bucketName string, dogRepository DogRepository) DogPhotoRepository {
	return &dogPhotoRepository{bucketName: bucketName, dogRepository: dogRepository}
}

func (r *dogPhotoRepository) Upload(id string, fileData io.Reader, contentType string) error {
	sess, err := common.CreateS3Session()
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
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
