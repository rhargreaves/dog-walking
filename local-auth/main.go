package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
)

type AuthorisedClaims struct {
	User   string
	Email  string
	Groups []string
}

func handleRequest(ctx context.Context,
	event events.APIGatewayV2CustomAuthorizerV1Request) (events.APIGatewayV2CustomAuthorizerIAMPolicyResponse, error) {
	if event.MethodArn == "" {
		return errorResponse("no MethodArn provided", ""), nil
	}

	authorisedClaims, err := authorise(event.AuthorizationToken, os.Getenv("LOCAL_JWT_SECRET"))
	if err != nil {
		return errorResponse(err.Error(), event.MethodArn), nil
	}

	return authorisedResponse(authorisedClaims.User,
		event.MethodArn,
		authorisedClaims.Email,
		authorisedClaims.Groups), nil
}

func authorise(authorisationToken string, jwtSecret string) (AuthorisedClaims, error) {
	if jwtSecret == "" {
		return AuthorisedClaims{}, errors.New("JWT secret missing")
	}
	if authorisationToken == "" {
		return AuthorisedClaims{}, errors.New("no AuthorizationToken provided")
	}
	tokenString := strings.TrimPrefix(authorisationToken, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return AuthorisedClaims{}, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return AuthorisedClaims{}, errors.New("token is invalid")
	}

	fmt.Println("token OK: claims:", claims)

	sub, ok := claims["sub"].(string)
	if !ok {
		return AuthorisedClaims{}, errors.New("token has no sub claim")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return AuthorisedClaims{}, errors.New("token has no email claim")
	}

	groups, ok := claims["cognito:groups"].([]any)
	if !ok {
		return AuthorisedClaims{}, errors.New("token has no cognito:groups claim")
	}

	return AuthorisedClaims{
		User:   sub,
		Email:  email,
		Groups: convertToStringSlice(groups),
	}, nil
}

func convertToStringSlice(slice []any) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = v.(string)
	}
	return result
}

func main() {
	lambda.Start(handleRequest)
}
