package breed_detector

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	"github.com/stretchr/testify/require"
)

const dummyImage = "dummy.jpeg"
const dummyBucket = "test-bucket"

var noSpecificDogBreedLabels = []*rekognition.Label{
	{
		Name:       aws.String("Dog"),
		Confidence: aws.Float64(90.0),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
}

type testRekognitionClient struct {
	rekognitioniface.RekognitionAPI
}

func NewTestRekognitionClient() rekognitioniface.RekognitionAPI {
	return &testRekognitionClient{}
}

func (m *testRekognitionClient) DetectLabels(input *rekognition.DetectLabelsInput) (*rekognition.DetectLabelsOutput, error) {
	return &rekognition.DetectLabelsOutput{
		Labels:            noSpecificDogBreedLabels,
		LabelModelVersion: aws.String("3.0"),
	}, nil
}

func (m *testRekognitionClient) DetectLabelsWithContext(ctx aws.Context, input *rekognition.DetectLabelsInput, opts ...request.Option) (*rekognition.DetectLabelsOutput, error) {
	return m.DetectLabels(input)
}

func TestDetectBreed_ReturnsNoSpecificBreedWhenOnlyDogLabelIsPresent(t *testing.T) {
	detector := NewBreedDetector(BreedDetectorConfig{BucketName: dummyBucket}, NewTestRekognitionClient())
	_, err := detector.DetectBreed(dummyImage)
	require.Equal(t, ErrNoSpecificBreedDetected, err)
}
