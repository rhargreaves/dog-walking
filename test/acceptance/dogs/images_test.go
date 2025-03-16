package dogs

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rhargreaves/dog-walking/test/acceptance/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestUploadImage(t *testing.T) {
	dog := createDog(t, "Rover")

	file, err := os.Open("../resources/dog.jpg")
	require.NoError(t, err)
	defer file.Close()

	url := fmt.Sprintf("%s/dogs/%s/photo", common.BaseUrl(), dog.ID)
	req, err := http.NewRequest("PUT", url, file)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "image/jpeg")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	t.Log("Image uploaded successfully")

	sess, err := createS3Session()
	require.NoError(t, err)

	s3Client := s3.New(sess)
	key := dog.ID
	bucket := os.Getenv("DOG_IMAGES_BUCKET")
	_, err = s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	assert.NoError(t, err, "Image should exist in S3 bucket: "+bucket+" at key: "+key)
}
