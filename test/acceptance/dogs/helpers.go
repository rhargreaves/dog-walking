package dogs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhargreaves/dog-walking/test/acceptance/common"
)

type Dog struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func createDog(t *testing.T, name string) Dog {
	resp := common.PostJson(t, "/dogs", Dog{Name: name})
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusCreated)

	var dog Dog
	common.DecodeJSON(t, resp, &dog)

	assert.Equal(t, name, dog.Name, "Expected dog name to be returned")
	assert.NotEmpty(t, dog.ID, "Expected dog ID to be returned")

	return dog
}

func uploadFile(t *testing.T, url string, file *os.File, contentType string) {
	req, err := http.NewRequest("PUT", url, file)
	require.NoError(t, err)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
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

func requireS3ObjectExists(t *testing.T, bucket string, key string) {
	sess, err := createS3Session()
	require.NoError(t, err)

	s3Client := s3.New(sess)
	_, err = s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	require.NoError(t, err,
		fmt.Sprintf("Object should exist in S3 bucket: %s key: %s", bucket, key))
}

func getS3Object(t *testing.T, bucket string, key string) []byte {
	sess, err := createS3Session()
	require.NoError(t, err)

	s3Client := s3.New(sess)
	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	require.NoError(t, err)
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	return body
}
