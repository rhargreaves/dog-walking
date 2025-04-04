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
)

type stubRekognitionClient struct {
	rekognitioniface.RekognitionAPI
	s3Svc *s3.S3
}

func NewStubRekognitionClient(s3Svc *s3.S3) rekognitioniface.RekognitionAPI {
	return &stubRekognitionClient{s3Svc: s3Svc}
}

func (m *stubRekognitionClient) DetectLabels(input *rekognition.DetectLabelsInput) (*rekognition.DetectLabelsOutput, error) {
	image, err := m.s3Svc.GetObject(&s3.GetObjectInput{
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
