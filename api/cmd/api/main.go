package main

// @title Dog Walking Service API
// @version 1.0
// @description API for managing dogs, etc
// @BasePath /
// @schemes https

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
	_ "github.com/rhargreaves/dog-walking/api/docs"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs"
	"github.com/rhargreaves/dog-walking/api/internal/rekognition_stub"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var ginLambda *ginadapter.GinLambdaV2

func init() {
	if !common.IsLocal() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(common.ErrorMiddleware)

	dogRepository := dogs.NewDynamoDBDogRepository(dogs.DynamoDBDogRepositoryConfig{
		TableName: os.Getenv("DOGS_TABLE_NAME"),
	})

	dogHandler := dogs.NewDogHandler(dogs.DogHandlerConfig{
		ImagesCdnBaseUrl: os.Getenv("CLOUDFRONT_BASE_URL"),
	}, dogRepository)

	dogPhotoUploader := dogs.NewDogPhotoUploader(dogs.S3PhotoUploaderConfig{
		BucketName: os.Getenv("DOG_IMAGES_BUCKET"),
	}, dogRepository)

	breedDetector := dogs.NewBreedDetector(dogs.BreedDetectorConfig{
		BucketName: os.Getenv("DOG_IMAGES_BUCKET"),
	}, newRekognitionClient())

	dogPhotoHandler := dogs.NewDogPhotoHandler(dogRepository, dogPhotoUploader, breedDetector)

	r.GET("/ping", pingHandler)
	r.GET("/dogs", dogHandler.ListDogs)
	r.GET("/dogs/:id", dogHandler.GetDog)
	r.POST("/dogs", dogHandler.CreateDog)
	r.PUT("/dogs/:id", dogHandler.UpdateDog)
	r.DELETE("/dogs/:id", dogHandler.DeleteDog)
	r.PUT("/dogs/:id/photo", dogPhotoHandler.UploadDogPhoto)
	r.POST("/dogs/:id/photo/detect-breed", dogPhotoHandler.DetectBreed)
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/api-docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api-docs/index.html")
	})

	ginLambda = ginadapter.NewV2(r)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}

// PingHandler godoc
// @Summary Health check endpoint
// @Description Returns OK if the API is running
// @Tags health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /ping [get]
func pingHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func newRekognitionClient() rekognitioniface.RekognitionAPI {
	if common.IsLocal() {
		return rekognition_stub.NewStubRekognitionClient()
	} else {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		}))
		return rekognition.New(sess)
	}
}
