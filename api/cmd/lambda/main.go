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
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/rhargreaves/dog-walking/api/docs"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/health"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var ginLambda *ginadapter.GinLambdaV2

func init() {
	isLocal := os.Getenv("USE_LOCALSTACK") == "true"
	dogHandler, dogPhotoHandler := createHandlers(isLocal)

	if !isLocal {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(common.ErrorMiddleware)

	r.GET("/ping", health.PingHandler)
	r.GET("/dogs", dogHandler.ListDogs)
	r.GET("/dogs/:id", dogHandler.GetDog)
	r.POST("/dogs", dogHandler.CreateDog)
	r.PUT("/dogs/:id", dogHandler.UpdateDog)
	r.DELETE("/dogs/:id", dogHandler.DeleteDog)
	r.PUT("/dogs/:id/photo", dogPhotoHandler.UploadDogPhoto)
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
