package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/rhargreaves/dog-walking/api/internal/common"
	"github.com/rhargreaves/dog-walking/api/internal/dogs"
)

var ginLambda *ginadapter.GinLambdaV2

func init() {
	r := gin.Default()
	r.Use(common.ErrorMiddleware)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	r.GET("/dogs", dogs.ListDogs)
	r.POST("/dogs", dogs.PostDog)

	ginLambda = ginadapter.NewV2(r)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
