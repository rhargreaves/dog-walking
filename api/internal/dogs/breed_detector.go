package dogs

import (
	"errors"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
)

var ErrNoDogDetected = errors.New("no dog detected")
var ErrNoSpecificBreedDetected = errors.New("no specific breed detected")

type BreedDetector interface {
	DetectBreed(id string) (*domain.BreedDetectionResult, error)
}

type BreedDetectorConfig struct {
	BucketName string
}

type breedDetector struct {
	config    *BreedDetectorConfig
	rekClient rekognitioniface.RekognitionAPI
}

func NewBreedDetector(breedDetectorConfig BreedDetectorConfig, rekClient rekognitioniface.RekognitionAPI) BreedDetector {
	return &breedDetector{
		config:    &breedDetectorConfig,
		rekClient: rekClient,
	}
}

func (d *breedDetector) DetectBreed(id string) (*domain.BreedDetectionResult, error) {
	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(d.config.BucketName),
				Name:   aws.String(id),
			},
		},
		MaxLabels:     aws.Int64(10),
		MinConfidence: aws.Float64(55.0),
		Settings: &rekognition.DetectLabelsSettings{
			GeneralLabels: &rekognition.GeneralLabelsSettings{
				LabelCategoryInclusionFilters: []*string{
					aws.String("Animals and Pets"),
				},
				LabelExclusionFilters: []*string{
					aws.String("Pet"), // ignore generic labels
					aws.String("Mammal"),
					aws.String("Canine"),
					aws.String("Animal"),
				},
			},
		},
	}

	result, err := d.rekClient.DetectLabels(input)
	if err != nil {
		return nil, fmt.Errorf("failed to detect labels: %w", err)
	}

	// First check if it's a dog
	isDog := false
	for _, label := range result.Labels {
		if *label.Name == "Dog" {
			isDog = true
			break
		}
	}
	if !isDog {
		return nil, ErrNoDogDetected
	}

	// Look for breed-specific labels
	var breedLabels []*rekognition.Label
	for _, label := range result.Labels {
		if len(label.Parents) > 0 {
			for _, parent := range label.Parents {
				if *parent.Name == "Dog" {
					breedLabels = append(breedLabels, label)
				}
			}
		}
	}

	if len(breedLabels) == 0 {
		return nil, ErrNoSpecificBreedDetected
	}

	// Sort by confidence (highest first)
	sort.Slice(breedLabels, func(i, j int) bool {
		return *breedLabels[i].Confidence > *breedLabels[j].Confidence
	})

	// Return the highest confidence breed
	return &domain.BreedDetectionResult{
		Breed:      *breedLabels[0].Name,
		Confidence: *breedLabels[0].Confidence,
	}, nil
}
