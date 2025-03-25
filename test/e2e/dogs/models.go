package dogs

type Socialization struct {
	GoodWithChildren  bool `json:"goodWithChildren"`
	GoodWithPuppies   bool `json:"goodWithPuppies"`
	GoodWithLargeDogs bool `json:"goodWithLargeDogs"`
	GoodWithSmallDogs bool `json:"goodWithSmallDogs"`
}

type CreateDogRequest struct {
	Name                string        `json:"name"`
	Breed               string        `json:"breed"`
	Sex                 string        `json:"sex" binding:"required,oneof=male female"`
	IsNeutered          bool          `json:"isNeutered"`
	EnergyLevel         int           `json:"energyLevel" binding:"required,min=1,max=5"`
	Size                string        `json:"size" binding:"required,oneof=small medium large"`
	Socialization       Socialization `json:"socialization"`
	SpecialInstructions string        `json:"specialInstructions,omitempty"`
	DateOfBirth         string        `json:"dateOfBirth,omitempty"`
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
	PhotoUrl            string        `json:"photoUrl"`
	PhotoHash           string        `json:"photoHash"`
}

type DogListResponse struct {
	Dogs      []DogResponse `json:"dogs"`
	NextToken string        `json:"nextToken"`
}

type DetectBreedRequest struct {
}

type DetectBreedResponse struct {
	Breed      string  `json:"breed"`
	Confidence float64 `json:"confidence"`
}
