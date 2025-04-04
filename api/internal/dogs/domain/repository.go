package domain

import "errors"

type DogList struct {
	Dogs      []Dog
	NextToken string
}

type Dog struct {
	ID                  string
	Name                string
	Breed               string
	PhotoHash           string
	Sex                 string
	IsNeutered          bool
	EnergyLevel         int
	Size                string
	Socialization       Socialization
	SpecialInstructions string
	DateOfBirth         string
	PhotoStatus         string
}

type Socialization struct {
	GoodWithChildren  bool
	GoodWithPuppies   bool
	GoodWithLargeDogs bool
	GoodWithSmallDogs bool
}

type DogRepository interface {
	Create(dog Dog) (*Dog, error)
	List(limit int, name string, nextToken string) (*DogList, error)
	Get(id string) (*Dog, error)
	Update(id string, dog *Dog) error
	Delete(id string) error
	UpdatePhotoHash(id string, photoHash string) error
	UpdatePhotoStatus(id string, photoStatus string) error
}

var ErrDogNotFound = errors.New("dog not found")
