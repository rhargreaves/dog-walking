package dogs

type Dog struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Breed     string `json:"breed"`
	PhotoUrl  string `json:"photoUrl"`
	PhotoHash string `json:"photoHash"`
}

type DogList struct {
	Dogs      []Dog  `json:"dogs"`
	NextToken string `json:"nextToken"`
}
