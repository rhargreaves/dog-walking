package dogs

type DogResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Breed     string `json:"breed"`
	PhotoUrl  string `json:"photoUrl"`
	PhotoHash string `json:"photoHash"`
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
