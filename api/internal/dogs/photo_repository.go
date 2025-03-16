package dogs

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func createS3Session() (*session.Session, error) {
	useLocalStack := os.Getenv("USE_LOCALSTACK") == "true"
	region := os.Getenv("AWS_REGION")
	if useLocalStack {
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

func (r *dogPhotoRepository) Upload(id string, fileData io.Reader, contentType string) error {
	sess, err := createS3Session()
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
