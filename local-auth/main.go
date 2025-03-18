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

func authorisedResponse(principalID string, methodArn string, email string, groups []string) events.APIGatewayV2CustomAuthorizerIAMPolicyResponse {
	return events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{
		PrincipalID: principalID,
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
			"userId": principalID,
			"email":  email,
			"groups": groups,
		},
	}
}

func errorResponse(errorMessage string, methodArn string) events.APIGatewayV2CustomAuthorizerIAMPolicyResponse {
	return events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{
		PrincipalID: "",
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
			"error": errorMessage,
		},
	}
}

func handleRequest(ctx context.Context,
	event events.APIGatewayV2CustomAuthorizerV1Request) (events.APIGatewayV2CustomAuthorizerIAMPolicyResponse, error) {
	fmt.Println("event:", event)
	methodArn := event.MethodArn
	fmt.Println("methodArn:", methodArn)

	jwtSecret := os.Getenv("LOCAL_JWT_SECRET")
	fmt.Println("local_jwt_secret: ", jwtSecret)

	tokenString := event.AuthorizationToken
	if tokenString == "" {
		return errorResponse("No AuthorizationToken provided", methodArn), nil
	}

	fmt.Println("raw tokenString:", tokenString)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("trimmed tokenString:", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		fmt.Println("error:", err)
		return errorResponse(err.Error(), methodArn), nil
	}

	fmt.Println("token:", token)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("token OK: claims:", claims)
		return authorisedResponse(
			claims["sub"].(string),
			methodArn,
			claims["email"].(string),
			convertToStringSlice(claims["cognito:groups"].([]interface{})),
		), nil
	}

	fmt.Println("token is invalid")
	return errorResponse("invalid token", methodArn), nil
}

func convertToStringSlice(slice []interface{}) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = v.(string)
	}
	return result
}

func main() {
	lambda.Start(handleRequest)
}
