package main

// @title Dog Walking Service API
// @version 1.0
// @description API for managing dogs, etc
// @BasePath /
// @schemes https

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
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
	isLocal := os.Getenv("USE_LOCALSTACK") == "true"
	region := mustGetenv("AWS_REGION")
	session, err := createSession(isLocal, false, region)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	dogRepository := dogs.NewDynamoDBDogRepository(dogs.DynamoDBDogRepositoryConfig{
		TableName: mustGetenv("DOGS_TABLE_NAME"),
	}, session)

	dogHandler := dogs.NewDogHandler(dogs.DogHandlerConfig{
		ImagesCdnBaseUrl: mustGetenv("CLOUDFRONT_BASE_URL"),
	}, dogRepository)

	s3session, err := createSession(isLocal, true, region)
	if err != nil {
		log.Fatalf("Failed to create S3 session: %v", err)
	}

	dogImagesBucket := mustGetenv("DOG_IMAGES_BUCKET")
	dogPhotoUploader := dogs.NewDogPhotoUploader(dogs.S3PhotoUploaderConfig{
		BucketName: dogImagesBucket,
	}, dogRepository, s3session)

	rekognitionClient := newRekognitionClient(isLocal, session, s3session)
	breedDetector := dogs.NewBreedDetector(dogs.BreedDetectorConfig{
		BucketName: dogImagesBucket,
	}, rekognitionClient)

	dogPhotoHandler := dogs.NewDogPhotoHandler(dogRepository, dogPhotoUploader, breedDetector)

	r := configureGin(isLocal, dogHandler, dogPhotoHandler)
	ginLambda = ginadapter.NewV2(r)
}

func configureGin(isLocal bool, dogHandler dogs.DogHandler, dogPhotoHandler dogs.DogPhotoHandler) *gin.Engine {
	if !isLocal {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(common.ErrorMiddleware)

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

	return r
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

func newRekognitionClient(isLocal bool, session *session.Session, s3session *session.Session) rekognitioniface.RekognitionAPI {
	if isLocal {
		return rekognition_stub.NewStubRekognitionClient(s3session)
	}
	return rekognition.New(session)
}

func createSession(isLocal bool, isS3 bool, region string) (*session.Session, error) {
	config := &aws.Config{
		Region: &region,
	}
	if !isLocal {
		return session.NewSession(config)
	}

	var endpoint string
	if isS3 {
		endpoint = mustGetenv("AWS_S3_ENDPOINT_URL")
	} else {
		endpoint = mustGetenv("AWS_ENDPOINT_URL")
	}
	config.Endpoint = aws.String(endpoint)
	config.Credentials = credentials.NewStaticCredentials("test", "test", "")
	config.DisableSSL = aws.Bool(true)
	return session.NewSession(config)
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Required environment variable not set: " + key)
	}
	return val
}
