package rekognition_stub

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

var ImageHashes = map[string]*[]*rekognition.Label{
	"443d6817146340599232418cfe7ef31b": &MrPeanutbutterLabels,
	"444514da01e4cd14434f4ede60d7998c": &HuskyLabels,
	"f1334bbc42d8894d475ccdcb154c8829": &NotAnAnimalLabels,
	"1a4d14f8b9b8233cadf7d24034a716fd": &CatLabels,
	"b1d862b24c5fe5ee2a06a650667a81a9": &DogWithGunLabels,
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

var MrPeanutbutterLabels = []*rekognition.Label{
	{
		Name:       aws.String("Airedale"),
		Confidence: aws.Float64(55.59829330444336),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Dog")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
			{Name: aws.String("Terrier")},
		},
	},
	{
		Name:       aws.String("Dog"),
		Confidence: aws.Float64(55.59829330444336),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
	{
		Name:       aws.String("Terrier"),
		Confidence: aws.Float64(55.59829330444336),
		Parents: []*rekognition.Parent{
			{Name: aws.String("Animal")},
			{Name: aws.String("Canine")},
			{Name: aws.String("Dog")},
			{Name: aws.String("Mammal")},
			{Name: aws.String("Pet")},
		},
	},
}

var DogWithGunLabels = []*rekognition.Label{
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
		Name:       aws.String("German Shepherd"),
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

var NotAnAnimalLabels = []*rekognition.Label{}
