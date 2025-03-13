package main

import (
	"context"
	"fmt"
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

	r.Use(func(c *gin.Context) {
		fmt.Printf("Request path: %s, method: %s\n", c.Request.URL.Path, c.Request.Method)
		c.Next()
	})

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			Message: "Hello World from the Dog Walking API!",
		})
	})

	r.POST("/dogs", func(c *gin.Context) {
		c.JSON(http.StatusOK, Dog{})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Route not found",
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})
	})

	ginLambda = ginadapter.New(r)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
