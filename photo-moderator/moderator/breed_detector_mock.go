package moderator

import "github.com/rhargreaves/dog-walking/photo-moderator/domain"

type TestBreedDetector struct {
	BreedDetector
	DetectBreedFunc func(id string) (*domain.BreedDetectionResult, error)
}

func (t *TestBreedDetector) DetectBreed(id string) (*domain.BreedDetectionResult, error) {
	return t.DetectBreedFunc(id)
}
