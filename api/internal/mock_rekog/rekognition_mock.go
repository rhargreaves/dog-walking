package mock_rekog

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
)

var mrPeanutbutterLabels = []*rekognition.Label{
	{
		Name:       aws.String("Airedale"),
		Confidence: aws.Float64(55.59829330444336),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Dog")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
			{Name: aws.String("Terrier")},
		},
	},
	{
		Name:       aws.String("Dog"),
		Confidence: aws.Float64(55.59829330444336),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
	{
		Name:       aws.String("Terrier"),
		Confidence: aws.Float64(55.59829330444336),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Dog")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
}

type mockRekognitionClient struct {
	rekognitioniface.RekognitionAPI
	labelsToOutput *[]*rekognition.Label
}

func NewMockRekognitionClient() rekognitioniface.RekognitionAPI {
	return &mockRekognitionClient{
		labelsToOutput: &mrPeanutbutterLabels,
	}
}

func NewMockRekognitionClientWithLabels(labels *[]*rekognition.Label) rekognitioniface.RekognitionAPI {
	return &mockRekognitionClient{
		labelsToOutput: labels,
	}
}

func (m *mockRekognitionClient) DetectLabels(input *rekognition.DetectLabelsInput) (*rekognition.DetectLabelsOutput, error) {
	return &rekognition.DetectLabelsOutput{
		Labels:            *m.labelsToOutput,
		LabelModelVersion: aws.String("3.0"),
	}, nil
}

func (m *mockRekognitionClient) DetectLabelsWithContext(ctx aws.Context, input *rekognition.DetectLabelsInput, opts ...request.Option) (*rekognition.DetectLabelsOutput, error) {
	return m.DetectLabels(input)
}
