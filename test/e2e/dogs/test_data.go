package dogs

var testDog = CreateOrUpdateDogRequest{
	Name:        "Rover",
	Breed:       "Labrador",
	Sex:         "male",
	IsNeutered:  true,
	EnergyLevel: 3,
	Size:        "medium",
	Socialization: Socialization{
		GoodWithChildren:  true,
		GoodWithPuppies:   true,
		GoodWithLargeDogs: true,
		GoodWithSmallDogs: true,
	},
	SpecialInstructions: "None",
	DateOfBirth:         "2020-01-01",
}

var testDog2 = CreateOrUpdateDogRequest{
	Name:        "Echo",
	Breed:       "Husky",
	Sex:         "male",
	IsNeutered:  false,
	EnergyLevel: 5,
	Size:        "large",
	Socialization: Socialization{
		GoodWithChildren:  false,
		GoodWithPuppies:   false,
		GoodWithLargeDogs: false,
		GoodWithSmallDogs: false,
	},
	SpecialInstructions: "Don't let him out of the house",
	DateOfBirth:         "2020-01-01",
}

var testListDog = CreateOrUpdateDogRequest{
	Name:        "ListTest",
	Breed:       "Husky",
	Sex:         "male",
	IsNeutered:  true,
	EnergyLevel: 3,
	Size:        "medium",
	Socialization: Socialization{
		GoodWithChildren:  true,
		GoodWithPuppies:   true,
		GoodWithLargeDogs: true,
		GoodWithSmallDogs: true,
	},
	SpecialInstructions: "None",
	DateOfBirth:         "2020-01-01",
}

var testNameFilterDog = CreateOrUpdateDogRequest{
	Name:        "NameFilterTest",
	Breed:       "Husky",
	Sex:         "male",
	IsNeutered:  true,
	EnergyLevel: 3,
	Size:        "medium",
	Socialization: Socialization{
		GoodWithChildren:  true,
		GoodWithPuppies:   true,
		GoodWithLargeDogs: true,
		GoodWithSmallDogs: true,
	},
	SpecialInstructions: "None",
	DateOfBirth:         "2020-01-01",
}
