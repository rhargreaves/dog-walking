package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs"
	"github.com/rhargreaves/dog-walking/api/internal/rekognition_stub"
)

var ginLambda *ginadapter.GinLambdaV2

func makeRekClient() rekognitioniface.RekognitionAPI {
	if os.Getenv("USE_LOCALSTACK") == "true" {
		return rekognition_stub.NewStubRekognitionClient()
	} else {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		}))
		return rekognition.New(sess)
	}
}

func init() {
	r := gin.Default()
	r.Use(common.ErrorMiddleware)

	dogRepository := dogs.NewDogRepository(os.Getenv("DOGS_TABLE_NAME"))
	dogHandler := dogs.NewDogHandler(dogRepository)

	dogPhotoRepository := dogs.NewDogPhotoRepository(os.Getenv("DOG_IMAGES_BUCKET"))
	breedDetector := dogs.NewBreedDetector(os.Getenv("DOG_IMAGES_BUCKET"), makeRekClient())
	dogPhotoHandler := dogs.NewDogPhotoHandler(dogRepository, dogPhotoRepository, breedDetector)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	r.GET("/dogs", dogHandler.ListDogs)
	r.GET("/dogs/:id", dogHandler.GetDog)
	r.POST("/dogs", dogHandler.CreateDog)
	r.PUT("/dogs/:id", dogHandler.UpdateDog)
	r.DELETE("/dogs/:id", dogHandler.DeleteDog)
	r.PUT("/dogs/:id/photo", dogPhotoHandler.UploadDogPhoto)
	r.POST("/dogs/:id/photo/detect-breed", dogPhotoHandler.DetectBreed)

	ginLambda = ginadapter.NewV2(r)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
