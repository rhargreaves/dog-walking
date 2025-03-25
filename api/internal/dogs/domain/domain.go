package domain

type Socialization struct {
	GoodWithChildren  bool
	GoodWithPuppies   bool
	GoodWithLargeDogs bool
	GoodWithSmallDogs bool
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
}

type DogList struct {
	Dogs      []Dog
	NextToken string
}

type BreedDetectionResult struct {
	Breed      string
	Confidence float64
}
