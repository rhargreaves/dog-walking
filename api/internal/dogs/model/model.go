package model

import (
	"fmt"

	"github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
)

type Socialization struct {
	GoodWithChildren  bool `json:"goodWithChildren"`
	GoodWithPuppies   bool `json:"goodWithPuppies"`
	GoodWithLargeDogs bool `json:"goodWithLargeDogs"`
	GoodWithSmallDogs bool `json:"goodWithSmallDogs"`
}

type CreateOrUpdateDogRequest struct {
	Name                string        `json:"name"`
	Breed               string        `json:"breed"`
	Sex                 string        `json:"sex" binding:"omitempty,oneof=male female"`
	IsNeutered          bool          `json:"isNeutered"`
	EnergyLevel         int           `json:"energyLevel" binding:"omitempty,min=1,max=5"`
	Size                string        `json:"size" binding:"omitempty,oneof=small medium large"`
	Socialization       Socialization `json:"socialization"`
	SpecialInstructions string        `json:"specialInstructions,omitempty"`
	DateOfBirth         string        `json:"dateOfBirth,omitempty"`
}

type DogListQuery struct {
	Limit     int    `form:"limit" default:"25" binding:"min=1,max=25"`
	NextToken string `form:"nextToken"`
	Name      string `form:"name"`
}

type DogResponse struct {
	ID                  string        `json:"id"`
	Name                string        `json:"name"`
	Breed               string        `json:"breed"`
	Sex                 string        `json:"sex" binding:"required,oneof=male female"`
	IsNeutered          bool          `json:"isNeutered"`
	EnergyLevel         int           `json:"energyLevel" binding:"required,min=1,max=5"`
	Size                string        `json:"size" binding:"required,oneof=small medium large"`
	Socialization       Socialization `json:"socialization"`
	SpecialInstructions string        `json:"specialInstructions,omitempty"`
	DateOfBirth         string        `json:"dateOfBirth,omitempty"`
	PhotoUrl            string        `json:"photoUrl,omitempty"`
	PhotoHash           string        `json:"photoHash,omitempty"`
	PhotoStatus         string        `json:"photoStatus,omitempty"`
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
	photoUrl := ""
	if dog.PhotoHash != "" {
		photoUrl = fmt.Sprintf("%s/%s?h=%s", imagesCdnBaseUrl, dog.ID, dog.PhotoHash)
	}
	return &DogResponse{
		ID:          dog.ID,
		Name:        dog.Name,
		Breed:       dog.Breed,
		Sex:         dog.Sex,
		IsNeutered:  dog.IsNeutered,
		EnergyLevel: dog.EnergyLevel,
		Size:        dog.Size,
		Socialization: Socialization{
			GoodWithChildren:  dog.Socialization.GoodWithChildren,
			GoodWithPuppies:   dog.Socialization.GoodWithPuppies,
			GoodWithLargeDogs: dog.Socialization.GoodWithLargeDogs,
			GoodWithSmallDogs: dog.Socialization.GoodWithSmallDogs,
		},
		SpecialInstructions: dog.SpecialInstructions,
		DateOfBirth:         dog.DateOfBirth,
		PhotoUrl:            photoUrl,
		PhotoHash:           dog.PhotoHash,
		PhotoStatus:         dog.PhotoStatus,
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

func FromCreateOrUpdateDogRequest(request *CreateOrUpdateDogRequest) *domain.Dog {
	return &domain.Dog{
		Name:        request.Name,
		Breed:       request.Breed,
		Sex:         request.Sex,
		IsNeutered:  request.IsNeutered,
		EnergyLevel: request.EnergyLevel,
		Size:        request.Size,
		Socialization: domain.Socialization{
			GoodWithChildren:  request.Socialization.GoodWithChildren,
			GoodWithPuppies:   request.Socialization.GoodWithPuppies,
			GoodWithLargeDogs: request.Socialization.GoodWithLargeDogs,
			GoodWithSmallDogs: request.Socialization.GoodWithSmallDogs,
		},
		SpecialInstructions: request.SpecialInstructions,
		DateOfBirth:         request.DateOfBirth,
	}
}
