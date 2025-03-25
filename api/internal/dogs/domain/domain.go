package domain

type Dog struct {
	ID        string
	Name      string
	Breed     string
	PhotoHash string
}

type DogList struct {
	Dogs      []Dog
	NextToken string
}

type BreedDetectionResult struct {
	Breed      string
	Confidence float64
}
