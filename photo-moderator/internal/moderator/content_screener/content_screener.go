package content_screener

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
)

type ContentScreener interface {
	ScreenImage(id string) (*ContentScreenerResult, error)
}

type ContentScreenerConfig struct {
	BucketName string
}

type contentScreener struct {
	config    *ContentScreenerConfig
	rekClient rekognitioniface.RekognitionAPI
}

type ContentScreenerResult struct {
	IsSafe bool
	Reason string
}

func NewContentScreener(config ContentScreenerConfig, rekClient rekognitioniface.RekognitionAPI) ContentScreener {
	return &contentScreener{
		config:    &config,
		rekClient: rekClient,
	}
}

func (c *contentScreener) ScreenImage(id string) (*ContentScreenerResult, error) {
	rekClient := c.rekClient

	resp, err := rekClient.DetectModerationLabels(&rekognition.DetectModerationLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(c.config.BucketName),
				Name:   aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to screen image for inappropriate content: %w", err)
	}

	if len(resp.ModerationLabels) == 0 {
		return &ContentScreenerResult{
			IsSafe: true,
			Reason: "",
		}, nil
	}

	return &ContentScreenerResult{
		IsSafe: false,
		Reason: *resp.ModerationLabels[0].Name,
	}, nil
}
