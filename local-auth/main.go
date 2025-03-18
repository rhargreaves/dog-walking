package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
)

func handleRequest(ctx context.Context,
	event events.APIGatewayV2CustomAuthorizerV1Request) (events.APIGatewayV2CustomAuthorizerIAMPolicyResponse, error) {
	fmt.Println("ctx:", ctx)
	fmt.Println("event:", event)
	methodArn := event.MethodArn
	fmt.Println("methodArn:", methodArn)

	jwtSecret := os.Getenv("LOCAL_JWT_SECRET")
	fmt.Println("local_jwt_secret: ", jwtSecret)

	tokenString := event.AuthorizationToken
	fmt.Println("raw tokenString:", tokenString)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("trimmed tokenString:", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		fmt.Println("error:", err)
		return events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{
			PrincipalID: "user",
			PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
				Version: "2012-10-17",
				Statement: []events.IAMPolicyStatement{
					{
						Action:   []string{"execute-api:Invoke"},
						Effect:   "Deny",
						Resource: []string{methodArn},
					},
				},
			},
			Context: map[string]interface{}{
				"error": err.Error(),
			},
		}, nil
	}

	fmt.Println("token:", token)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("token OK: claims:", claims)
		return events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{
			PrincipalID: claims["sub"].(string),
			PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
				Version: "2012-10-17",
				Statement: []events.IAMPolicyStatement{
					{
						Action:   []string{"execute-api:Invoke"},
						Effect:   "Allow",
						Resource: []string{methodArn},
					},
				},
			},
			Context: map[string]interface{}{
				"userId": claims["sub"],
				"email":  claims["email"],
				"groups": claims["cognito:groups"],
			},
		}, nil
	}

	fmt.Println("token is invalid")
	return events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Deny",
					Resource: []string{methodArn},
				},
			},
		},
		Context: map[string]interface{}{
			"error": "invalid token",
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
