package rekognition_stub

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

var ImageHashToModerationLabels = map[string]*[]*rekognition.ModerationLabel{
	"443d6817146340599232418cfe7ef31b": &NoModerationLabels,
	"444514da01e4cd14434f4ede60d7998c": &NoModerationLabels,
	"f1334bbc42d8894d475ccdcb154c8829": &NoModerationLabels,
	"1a4d14f8b9b8233cadf7d24034a716fd": &NoModerationLabels,
	"b1d862b24c5fe5ee2a06a650667a81a9": &DogWithGunModerationLabels,
}

var NoModerationLabels = []*rekognition.ModerationLabel{}

var DogWithGunModerationLabels = []*rekognition.ModerationLabel{
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
