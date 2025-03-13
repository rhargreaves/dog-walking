package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

type Dog struct {
}

var ginLambda *ginadapter.GinLambda

func init() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			Message: "Hello World from the Dog Walking API!",
		})
	})

	r.POST("/dogs", func(c *gin.Context) {
		c.JSON(http.StatusOK, Dog{})
	})

	ginLambda = ginadapter.New(r)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
