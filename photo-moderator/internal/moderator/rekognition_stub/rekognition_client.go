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

func (m *stubRekognitionClient) getImageHash(bucket string, key string) (string, error) {
	image, err := m.s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}
	imageBytes, err := io.ReadAll(image.Body)
	if err != nil {
		return "", err
	}
	hash := md5.Sum(imageBytes)
	hashString := hex.EncodeToString(hash[:])
	return hashString, nil
}

func (m *stubRekognitionClient) DetectLabels(input *rekognition.DetectLabelsInput) (*rekognition.DetectLabelsOutput, error) {
	hashString, err := m.getImageHash(*input.Image.S3Object.Bucket, *input.Image.S3Object.Name)
	if err != nil {
		return nil, err
	}

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

func (m *stubRekognitionClient) DetectModerationLabels(input *rekognition.DetectModerationLabelsInput) (*rekognition.DetectModerationLabelsOutput, error) {
	hashString, err := m.getImageHash(*input.Image.S3Object.Bucket, *input.Image.S3Object.Name)
	if err != nil {
		return nil, err
	}

	moderationLabels, ok := ImageHashToModerationLabels[hashString]
	if !ok {
		return nil, errors.New("Stub Rekognition client: moderation labels hash (" + hashString + ") not found")
	}

	return &rekognition.DetectModerationLabelsOutput{
		ModerationLabels:       *moderationLabels,
		ModerationModelVersion: aws.String("7.0"),
		ContentTypes:           []*rekognition.ContentType{},
	}, nil
}

func (m *stubRekognitionClient) DetectModerationLabelsWithContext(ctx aws.Context, input *rekognition.DetectModerationLabelsInput, opts ...request.Option) (*rekognition.DetectModerationLabelsOutput, error) {
	return m.DetectModerationLabels(input)
}
