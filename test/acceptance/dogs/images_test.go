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

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Endpoint:    aws.String("http://s3.localhost.localstack.cloud:4566"),
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
	})
	require.NoError(t, err)

	s3Client := s3.New(sess)
	key := fmt.Sprintf("%s", dog.ID)
	_, err = s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String("local-dog-images"),
		Key:    aws.String(key),
	})
	assert.NoError(t, err, "Image should exist in S3 bucket at key: "+key)
}
