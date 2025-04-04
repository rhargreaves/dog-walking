package rekognition_stub

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

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

var NotAnAnimalLabels = []*rekognition.Label{}
