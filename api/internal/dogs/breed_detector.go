package dogs

import (
	"errors"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
)

var ErrNoDogDetected = errors.New("no dog detected")
var ErrNoSpecificBreedDetected = errors.New("no specific breed detected")

type BreedDetector interface {
	DetectBreed(id string) (string, float64, error)
}

type breedDetector struct {
	bucketName string
	rekClient  rekognitioniface.RekognitionAPI
}

func NewBreedDetector(bucketName string, rekClient rekognitioniface.RekognitionAPI) BreedDetector {
	return &breedDetector{
		bucketName: bucketName,
		rekClient:  rekClient,
	}
}

func (d *breedDetector) DetectBreed(id string) (string, float64, error) {
	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(d.bucketName),
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
		return "", 0, fmt.Errorf("failed to detect labels: %w", err)
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
		return "", 0, ErrNoDogDetected
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
		return "", 0, ErrNoSpecificBreedDetected
	}

	// Sort by confidence (highest first)
	sort.Slice(breedLabels, func(i, j int) bool {
		return *breedLabels[i].Confidence > *breedLabels[j].Confidence
	})

	// Return the highest confidence breed
	topBreed := *breedLabels[0].Name
	confidence := *breedLabels[0].Confidence

	return topBreed, confidence, nil
}
