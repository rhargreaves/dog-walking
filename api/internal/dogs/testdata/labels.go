package testdata

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

var GermanShepherdLabels = []*rekognition.Label{
	{
		Name:       aws.String("Dog"),
		Confidence: aws.Float64(97.3577651977539),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
	{
		Name:       aws.String("German Shepherd"),
		Confidence: aws.Float64(80.85753631591797),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Dog")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
}

var HuskyLabels = []*rekognition.Label{
	{
		Name:       aws.String("Dog"),
		Confidence: aws.Float64(99.54544067382812),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
	{
		Name:       aws.String("Husky"),
		Confidence: aws.Float64(99.54544067382812),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Dog")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
}

var TerrierLabels = []*rekognition.Label{
	{
		Name:       aws.String("Dog"),
		Confidence: aws.Float64(98.3401870727539),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
	{
		Name:       aws.String("Terrier"),
		Confidence: aws.Float64(86.38780212402344),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Dog")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
	{
		Name:       aws.String("Airedale"),
		Confidence: aws.Float64(57.59425354003906),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Dog")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
			{Name: aws.String("Terrier")},
		},
	},
}

var CatLabels = []*rekognition.Label{
	{
		Name:       aws.String("Cat"),
		Confidence: aws.Float64(96.29752349853516),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
}

var NoSpecificDogBreedLabels = []*rekognition.Label{
	{
		Name:       aws.String("Dog"),
		Confidence: aws.Float64(90.0),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
}

var NotAnAnimal = []*rekognition.Label{}
