package rekognition_stub

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	"github.com/aws/aws-sdk-go/service/s3"
)

type stubRekognitionClient struct {
	rekognitioniface.RekognitionAPI
	labelsToOutput *[]*rekognition.Label
}

func NewStubRekognitionClient() rekognitioniface.RekognitionAPI {
	return &stubRekognitionClient{
		labelsToOutput: &MrPeanutbutterLabels,
	}
}

func NewStubRekognitionClientWithLabels(labels *[]*rekognition.Label) rekognitioniface.RekognitionAPI {
	return &stubRekognitionClient{
		labelsToOutput: labels,
	}
}

func (m *stubRekognitionClient) DetectLabels(input *rekognition.DetectLabelsInput) (*rekognition.DetectLabelsOutput, error) {
	session, err := createS3Session()
	if err != nil {
		return nil, err
	}
	s3client := s3.New(session)
	image, err := s3client.GetObject(&s3.GetObjectInput{
		Bucket: input.Image.S3Object.Bucket,
		Key:    input.Image.S3Object.Name,
	})
	if err != nil {
		return nil, err
	}

	imageBytes, err := io.ReadAll(image.Body)
	if err != nil {
		return nil, err
	}

	hash := md5.Sum(imageBytes)
	hashString := hex.EncodeToString(hash[:])
	labels, ok := ImageHashes[hashString]
	if !ok {
		return nil, errors.New("image hash (" + hashString + ") not found")
	}

	return &rekognition.DetectLabelsOutput{
		Labels:            *labels,
		LabelModelVersion: aws.String("3.0"),
	}, nil
}

func (m *stubRekognitionClient) DetectLabelsWithContext(ctx aws.Context, input *rekognition.DetectLabelsInput, opts ...request.Option) (*rekognition.DetectLabelsOutput, error) {
	return m.DetectLabels(input)
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
