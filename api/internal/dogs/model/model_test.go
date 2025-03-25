package model

import (
	"testing"

	"github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
	"github.com/stretchr/testify/assert"
)

const testImagesBaseUrl = "http://images.localhost"

func TestToDogResponse_AddsPhotoUrlIfHashIsPresent(t *testing.T) {
	dog := &domain.Dog{
		ID:        "1",
		Name:      "Fido",
		PhotoHash: "123",
	}

	response := ToDogResponse(dog, testImagesBaseUrl)

	assert.Equal(t, testImagesBaseUrl+"/1?h=123", response.PhotoUrl)
}

func TestToDogResponse_DoesNotAddPhotoUrlIfHashIsEmpty(t *testing.T) {
	dog := &domain.Dog{
		ID:        "1",
		Name:      "Fido",
		PhotoHash: "",
	}

	response := ToDogResponse(dog, testImagesBaseUrl)

	assert.Equal(t, "", response.PhotoUrl)
}

func TestToDogListResponse_AddsPhotoUrlIfHashIsPresent(t *testing.T) {
	dog := &domain.Dog{
		ID:        "1",
		Name:      "Fido",
		PhotoHash: "123",
	}
	dogList := &domain.DogList{
		Dogs: []domain.Dog{*dog},
	}

	response := ToDogListResponse(dogList, testImagesBaseUrl)

	assert.Equal(t, testImagesBaseUrl+"/1?h=123", response.Dogs[0].PhotoUrl)
}

func TestToDogListResponse_DoesNotAddPhotoUrlIfHashIsEmpty(t *testing.T) {
	dog := &domain.Dog{
		ID:        "1",
		Name:      "Fido",
		PhotoHash: "",
	}
	dogList := &domain.DogList{
		Dogs: []domain.Dog{*dog},
	}

	response := ToDogListResponse(dogList, testImagesBaseUrl)

	assert.Equal(t, "", response.Dogs[0].PhotoUrl)
}
