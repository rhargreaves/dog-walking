package dogs

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/testdata"
	"github.com/rhargreaves/dog-walking/api/internal/mock_rekog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dummyImage = "dummy.jpeg"
const dummyBucket = "test-bucket"

func TestDetectBreed(t *testing.T) {
	testCases := []struct {
		name          string
		labels        []*rekognition.Label
		expectedBreed string
		expectedConf  float64
		expectedErr   error
	}{
		{
			name:        "not an animal",
			labels:      testdata.NotAnAnimal,
			expectedErr: ErrNoDogDetected,
		},
		{
			name:        "cat",
			labels:      testdata.CatLabels,
			expectedErr: ErrNoDogDetected,
		},
		{
			name:          "german shepherd",
			labels:        testdata.GermanShepherdLabels,
			expectedBreed: "German Shepherd",
			expectedConf:  80.85753631591797,
		},
		{
			name:          "husky",
			labels:        testdata.HuskyLabels,
			expectedBreed: "Husky",
			expectedConf:  99.54544067382812,
		},
		{
			name:          "terrier",
			labels:        testdata.TerrierLabels,
			expectedBreed: "Terrier",
			expectedConf:  86.38780212402344,
		},
		{
			name:        "no specific breed",
			labels:      testdata.NoSpecificDogBreedLabels,
			expectedErr: ErrNoSpecificBreedDetected,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewBreedDetector(dummyBucket, mock_rekog.NewMockRekognitionClientWithLabels(&tt.labels))
			breed, conf, err := detector.DetectBreed(dummyImage)

			if tt.expectedErr != nil {
				require.Equal(t, tt.expectedErr, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedBreed, breed)
			assert.Equal(t, tt.expectedConf, conf)
		})
	}
}

func TestDetectBreed_ReturnsDefaultBreed(t *testing.T) {
	detector := NewBreedDetector(dummyBucket, mock_rekog.NewMockRekognitionClient())
	breed, confidence, err := detector.DetectBreed(dummyImage)
	assert.NoError(t, err)
	assert.Equal(t, "Airedale", breed)
	assert.Equal(t, 55.59829330444336, confidence)
}
