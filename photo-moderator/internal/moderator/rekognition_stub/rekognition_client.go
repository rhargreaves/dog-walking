package rekognition_stub

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

func (r *stubRekognitionClient) DetectLabels(input *rekognition.DetectLabelsInput) (*rekognition.DetectLabelsOutput, error) {
	imageClassification, err := r.getImageClassification(*input.Image.S3Object.Bucket, *input.Image.S3Object.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to detect labels: %w", err)
	}

	return &rekognition.DetectLabelsOutput{
		Labels:            imageClassification.Labels,
		LabelModelVersion: aws.String("3.0"),
	}, nil
}

func (r *stubRekognitionClient) DetectLabelsWithContext(ctx aws.Context, input *rekognition.DetectLabelsInput, opts ...request.Option) (*rekognition.DetectLabelsOutput, error) {
	return r.DetectLabels(input)
}

func (r *stubRekognitionClient) DetectModerationLabels(input *rekognition.DetectModerationLabelsInput) (*rekognition.DetectModerationLabelsOutput, error) {
	imageClassification, err := r.getImageClassification(*input.Image.S3Object.Bucket, *input.Image.S3Object.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to detect moderation labels: %w", err)
	}

	return &rekognition.DetectModerationLabelsOutput{
		ModerationLabels:       imageClassification.ModerationLabels,
		ModerationModelVersion: aws.String("7.0"),
		ContentTypes:           []*rekognition.ContentType{},
	}, nil
}

func (r *stubRekognitionClient) DetectModerationLabelsWithContext(ctx aws.Context, input *rekognition.DetectModerationLabelsInput, opts ...request.Option) (*rekognition.DetectModerationLabelsOutput, error) {
	return r.DetectModerationLabels(input)
}

func (r *stubRekognitionClient) getImageClassification(bucket string, key string) (*ImageClassification, error) {
	image, err := r.s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get image from S3: %w", err)
	}

	imageBytes, err := io.ReadAll(image.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image bytes: %w", err)
	}
	hash := md5.Sum(imageBytes)
	hashString := hex.EncodeToString(hash[:])
	imageClassification, ok := ImageClassifications[hashString]
	if !ok {
		return nil, fmt.Errorf("stubRekognitionClient: image classification for hash (%s) not found", hashString)
	}
	return imageClassification, nil
}
