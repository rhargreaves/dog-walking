package rekognition_stub

import (
	"github.com/aws/aws-sdk-go/service/rekognition"
)

var ImageHashes = map[string]*[]*rekognition.Label{
	"443d6817146340599232418cfe7ef31b": &MrPeanutbutterLabels,
	"444514da01e4cd14434f4ede60d7998c": &HuskyLabels,
	"f1334bbc42d8894d475ccdcb154c8829": &NotAnAnimalLabels,
	"1a4d14f8b9b8233cadf7d24034a716fd": &CatLabels,
}
