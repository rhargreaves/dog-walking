package models

type Dog struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Breed     string `json:"breed,omitempty"`
	PhotoUrl  string `json:"photoUrl,omitempty"`
	PhotoHash string `json:"photoHash,omitempty"`
}

type DogList struct {
	Dogs      []Dog  `json:"dogs"`
	NextToken string `json:"nextToken"`
}
