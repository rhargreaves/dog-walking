package mocks

import (
	"github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/breed_detector"
)

type MockBreedDetector struct {
	breed_detector.BreedDetector
	DetectBreedFunc func(id string) (*breed_detector.BreedDetectionResult, error)
}

func (t *MockBreedDetector) DetectBreed(id string) (*breed_detector.BreedDetectionResult, error) {
	return t.DetectBreedFunc(id)
}
