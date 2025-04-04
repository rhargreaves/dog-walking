package content_screener

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	"github.com/stretchr/testify/require"
)

type testRekognitionClient struct {
	rekognitioniface.RekognitionAPI
	ModerationLabels []*rekognition.ModerationLabel
}

func NewTestRekognitionClient() *testRekognitionClient {
	return &testRekognitionClient{}
}

func (m *testRekognitionClient) DetectModerationLabels(input *rekognition.DetectModerationLabelsInput) (*rekognition.DetectModerationLabelsOutput, error) {
	return &rekognition.DetectModerationLabelsOutput{
		ModerationLabels: m.ModerationLabels,
	}, nil
}

func (m *testRekognitionClient) ReturnsModerationLabels(labels []*rekognition.ModerationLabel) {
	m.ModerationLabels = labels
}

func TestScreenImage_AnyModerationLabels_ReturnsUnsafe(t *testing.T) {
	rekClient := NewTestRekognitionClient()
	rekClient.ReturnsModerationLabels([]*rekognition.ModerationLabel{
		{Name: aws.String("ModerationLabel")},
	})
	contentScreener := NewContentScreener(ContentScreenerConfig{BucketName: "test-bucket"}, rekClient)

	result, err := contentScreener.ScreenImage("test-image")
	require.NoError(t, err)
	require.False(t, result.IsSafe)
	require.Equal(t, "ModerationLabel", result.Reason)
}

func TestScreenImage_NoModerationLabels_ReturnsSafe(t *testing.T) {
	rekClient := NewTestRekognitionClient()
	rekClient.ReturnsModerationLabels([]*rekognition.ModerationLabel{})
	contentScreener := NewContentScreener(ContentScreenerConfig{BucketName: "test-bucket"}, rekClient)

	result, err := contentScreener.ScreenImage("test-image")
	require.NoError(t, err)
	require.True(t, result.IsSafe)
	require.Equal(t, "", result.Reason)
}
