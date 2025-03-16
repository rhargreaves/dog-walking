package mock_rekog

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
)

//go:embed mock_labels.json
var mockLabelsJSON []byte

type RekognitionLabelsResponse struct {
	Labels            []*rekognition.Label `json:"Labels"`
	LabelModelVersion string               `json:"LabelModelVersion"`
}

type mockRekognitionClient struct {
	rekognitioniface.RekognitionAPI
}

func NewMockRekognitionClient() rekognitioniface.RekognitionAPI {
	return &mockRekognitionClient{}
}

func (m *mockRekognitionClient) DetectLabels(input *rekognition.DetectLabelsInput) (*rekognition.DetectLabelsOutput, error) {
	var response RekognitionLabelsResponse
	if err := json.Unmarshal(mockLabelsJSON, &response); err != nil {
		return nil, fmt.Errorf("failed to parse mock labels JSON: %w", err)
	}

	return &rekognition.DetectLabelsOutput{
		Labels:            response.Labels,
		LabelModelVersion: &response.LabelModelVersion,
	}, nil
}

func (m *mockRekognitionClient) DetectLabelsWithContext(ctx aws.Context, input *rekognition.DetectLabelsInput, opts ...request.Option) (*rekognition.DetectLabelsOutput, error) {
	return m.DetectLabels(input)
}
