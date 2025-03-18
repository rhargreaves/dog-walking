package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const methodArn = "arn:aws:execute-api:us-east-1:123456789012:api-id/stage/method/resourcepath"

func TestHandleRequest(t *testing.T) {
	event := events.APIGatewayV2CustomAuthorizerV1Request{
		MethodArn:          methodArn,
		AuthorizationToken: createTestJWT(t),
	}

	response, err := handleRequest(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, "test-user-id", response.PrincipalID)
	assert.Equal(t, "Allow", response.PolicyDocument.Statement[0].Effect)
	assert.Equal(t, []string{"execute-api:Invoke"}, response.PolicyDocument.Statement[0].Action)
	assert.Equal(t, "test-user-id", response.Context["userId"])
	assert.Equal(t, "test@example.com", response.Context["email"])
	assert.Equal(t, []string{"Users"}, response.Context["groups"])
}

func TestHandleRequest_MissingToken(t *testing.T) {
	event := events.APIGatewayV2CustomAuthorizerV1Request{
		MethodArn: methodArn,
	}

	response, err := handleRequest(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, "", response.PrincipalID)
	assert.Equal(t, "Deny", response.PolicyDocument.Statement[0].Effect)
	assert.Equal(t, "no AuthorizationToken provided", response.Context["error"])
}
func TestHandleRequest_MissingMethodArn(t *testing.T) {
	event := events.APIGatewayV2CustomAuthorizerV1Request{}

	response, err := handleRequest(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, "", response.PrincipalID)
	assert.Equal(t, "Deny", response.PolicyDocument.Statement[0].Effect)
	assert.Equal(t, "no MethodArn provided", response.Context["error"])
}

func TestHandleRequest_InvalidJWT(t *testing.T) {
	event := events.APIGatewayV2CustomAuthorizerV1Request{
		MethodArn:          methodArn,
		AuthorizationToken: "Bearer foo",
	}

	response, err := handleRequest(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, "", response.PrincipalID)
	assert.Equal(t, "Deny", response.PolicyDocument.Statement[0].Effect)
	assert.Equal(t, "token is malformed: token contains an invalid number of segments", response.Context["error"])
}

func TestHandleRequest_TokenHasExpired(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * -24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("LOCAL_JWT_SECRET")))
	require.NoError(t, err)

	event := events.APIGatewayV2CustomAuthorizerV1Request{
		MethodArn:          methodArn,
		AuthorizationToken: tokenString,
	}

	response, err := handleRequest(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, "", response.PrincipalID)
	assert.Equal(t, "Deny", response.PolicyDocument.Statement[0].Effect)
	assert.Equal(t, "token has invalid claims: token is expired", response.Context["error"])
}

func createTestJWT(t *testing.T) string {
	claims := jwt.MapClaims{
		"sub":            "test-user-id",
		"email":          "test@example.com",
		"cognito:groups": []string{"Users"},
		"exp":            time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("LOCAL_JWT_SECRET")))
	require.NoError(t, err)
	return tokenString
}
