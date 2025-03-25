package model

import (
	"fmt"

	"github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
)

type DogRequest struct {
	Name  string `json:"name" binding:"required"`
	Breed string `json:"breed,omitempty"`
}

type DogResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Breed     string `json:"breed,omitempty"`
	PhotoUrl  string `json:"photoUrl,omitempty"`
	PhotoHash string `json:"photoHash,omitempty"`
}

type DogListResponse struct {
	Dogs      []DogResponse `json:"dogs"`
	NextToken string        `json:"nextToken"`
}

type BreedDetectionResultResponse struct {
	ID         string  `json:"id"`
	Breed      string  `json:"breed"`
	Confidence float64 `json:"confidence"`
}

func ToDogResponse(dog *domain.Dog, imagesCdnBaseUrl string) *DogResponse {
	photoUrl := fmt.Sprintf("%s/%s?h=%s",
		imagesCdnBaseUrl, dog.ID, dog.PhotoHash)
	return &DogResponse{
		ID:        dog.ID,
		Name:      dog.Name,
		Breed:     dog.Breed,
		PhotoUrl:  photoUrl,
		PhotoHash: dog.PhotoHash,
	}
}

func ToDogListResponse(dogs *domain.DogList, imagesCdnBaseUrl string) *DogListResponse {
	dogResponses := make([]DogResponse, len(dogs.Dogs))
	for i, dog := range dogs.Dogs {
		dogResponses[i] = *ToDogResponse(&dog, imagesCdnBaseUrl)
	}
	return &DogListResponse{
		Dogs:      dogResponses,
		NextToken: dogs.NextToken,
	}
}

func ToBreedDetectionResultResponse(id string, breedResult *domain.BreedDetectionResult) *BreedDetectionResultResponse {
	return &BreedDetectionResultResponse{
		ID:         id,
		Breed:      breedResult.Breed,
		Confidence: breedResult.Confidence,
	}
}
