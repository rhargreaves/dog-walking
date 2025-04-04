package mocks

import (
	"github.com/rhargreaves/dog-walking/photo-moderator/domain"
	"github.com/rhargreaves/dog-walking/photo-moderator/moderator/breed_detector"
)

type MockBreedDetector struct {
	breed_detector.BreedDetector
	DetectBreedFunc func(id string) (*domain.BreedDetectionResult, error)
}

func (t *MockBreedDetector) DetectBreed(id string) (*domain.BreedDetectionResult, error) {
	return t.DetectBreedFunc(id)
}
