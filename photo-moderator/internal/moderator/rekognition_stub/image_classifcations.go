package rekognition_stub

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type ImageClassification struct {
	Labels           []*rekognition.Label
	ModerationLabels []*rekognition.ModerationLabel
}

var ImageClassifications = map[string]*ImageClassification{
	"443d6817146340599232418cfe7ef31b": {
		Labels:           mrPeanutbutterLabels,
		ModerationLabels: noModerationLabels,
	},
	"444514da01e4cd14434f4ede60d7998c": {
		Labels:           huskyLabels,
		ModerationLabels: noModerationLabels,
	},
	"f1334bbc42d8894d475ccdcb154c8829": {
		Labels:           notAnAnimalLabels,
		ModerationLabels: noModerationLabels,
	},
	"1a4d14f8b9b8233cadf7d24034a716fd": {
		Labels:           catLabels,
		ModerationLabels: noModerationLabels,
	},
	"b1d862b24c5fe5ee2a06a650667a81a9": {
		Labels:           dogWithGunLabels,
		ModerationLabels: dogWithGunModerationLabels,
	},
}

var huskyLabels = []*rekognition.Label{
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

var catLabels = []*rekognition.Label{
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

var mrPeanutbutterLabels = []*rekognition.Label{
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

var dogWithGunLabels = []*rekognition.Label{
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

var notAnAnimalLabels = []*rekognition.Label{}

var dogWithGunModerationLabels = []*rekognition.ModerationLabel{
	{
		Confidence:    aws.Float64(99.43720245361328),
		Name:          aws.String("Weapons"),
		ParentName:    aws.String("Violence"),
		TaxonomyLevel: aws.Int64(2),
	},
	{
		Confidence:    aws.Float64(99.43720245361328),
		Name:          aws.String("Violence"),
		ParentName:    aws.String(""),
		TaxonomyLevel: aws.Int64(1),
	},
}

var noModerationLabels = []*rekognition.ModerationLabel{}
