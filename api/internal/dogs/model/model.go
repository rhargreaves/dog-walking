package model

import (
	"fmt"
	"os"

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

func ToDogResponse(dog *domain.Dog) *DogResponse {
	photoUrl := fmt.Sprintf("%s/%s?h=%s",
		os.Getenv("CLOUDFRONT_BASE_URL"), dog.ID, dog.PhotoHash)
	return &DogResponse{
		ID:        dog.ID,
		Name:      dog.Name,
		Breed:     dog.Breed,
		PhotoUrl:  photoUrl,
		PhotoHash: dog.PhotoHash,
	}
}

func ToDogListResponse(dogs *domain.DogList) *DogListResponse {
	dogResponses := make([]DogResponse, len(dogs.Dogs))
	for i, dog := range dogs.Dogs {
		dogResponses[i] = *ToDogResponse(&dog)
	}
	return &DogListResponse{
		Dogs:      dogResponses,
		NextToken: dogs.NextToken,
	}
}
