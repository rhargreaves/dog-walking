package rekognition_stub

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rhargreaves/dog-walking/api/internal/common"
)

type stubRekognitionClient struct {
	rekognitioniface.RekognitionAPI
}

func NewStubRekognitionClient() rekognitioniface.RekognitionAPI {
	return &stubRekognitionClient{}
}

func (m *stubRekognitionClient) DetectLabels(input *rekognition.DetectLabelsInput) (*rekognition.DetectLabelsOutput, error) {
	session, err := common.CreateS3Session()
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
		return nil, errors.New("Stub Rekognition client: image hash (" + hashString + ") not found")
	}

	return &rekognition.DetectLabelsOutput{
		Labels:            *labels,
		LabelModelVersion: aws.String("3.0"),
	}, nil
}

func (m *stubRekognitionClient) DetectLabelsWithContext(ctx aws.Context, input *rekognition.DetectLabelsInput, opts ...request.Option) (*rekognition.DetectLabelsOutput, error) {
	return m.DetectLabels(input)
}
