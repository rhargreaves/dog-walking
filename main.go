package main

import (
	"context"
	"encoding/json"
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
	fmt.Printf("API Gateway request: method=%s, path=%s, resource=%s\n",
		req.HTTPMethod, req.Path, req.Resource)
	fmt.Printf("API Gateway headers: %v\n", req.Headers)

	ctxJSON, _ := json.MarshalIndent(req.RequestContext, "", "  ")
	fmt.Printf("API Gateway request context: %s\n", string(ctxJSON))
	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	fmt.Printf("Full API Gateway request: %s\n", string(reqJSON))

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
