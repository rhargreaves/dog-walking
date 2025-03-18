package auth

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/require"
)

func GetCognitoJWT(t *testing.T) string {
	poolName := os.Getenv("COGNITO_USER_POOL_NAME")
	require.NotEmpty(t, poolName, "COGNITO_USER_POOL_NAME environment variable is required")
	clientName := os.Getenv("COGNITO_CLIENT_NAME")
	require.NotEmpty(t, clientName, "COGNITO_CLIENT_NAME environment variable is required")

	sess := session.Must(session.NewSession())
	cognito := cognitoidentityprovider.New(sess)
	userPoolId := findUserPoolByName(t, cognito, poolName)
	clientId := findClientByName(t, cognito, userPoolId, clientName)

	username := os.Getenv("COGNITO_USERNAME")
	require.NotEmpty(t, username, "COGNITO_USERNAME environment variable is required")
	password := os.Getenv("COGNITO_PASSWORD")
	require.NotEmpty(t, password, "COGNITO_PASSWORD environment variable is required")
	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String(clientId),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
	}

	authOutput, err := cognito.InitiateAuth(authInput)
	require.NoError(t, err, "failed to authenticate with Cognito")
	require.NotNil(t, authOutput.AuthenticationResult, "no authentication result received")
	return *authOutput.AuthenticationResult.IdToken
}

func findUserPoolByName(t *testing.T, cognito *cognitoidentityprovider.CognitoIdentityProvider,
	poolName string) string {
	listPoolsInput := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int64(1),
	}
	listPoolsOutput, err := cognito.ListUserPools(listPoolsInput)
	require.NoError(t, err, "failed to list user pools")

	var userPoolId string
	for _, pool := range listPoolsOutput.UserPools {
		if *pool.Name == poolName {
			userPoolId = *pool.Id
			break
		}
	}
	require.NotEmpty(t, userPoolId,
		fmt.Sprintf("could not find user pool with name '%s'", poolName))
	return userPoolId
}

func findClientByName(t *testing.T, cognito *cognitoidentityprovider.CognitoIdentityProvider,
	poolId string, clientName string) string {
	listClientsInput := &cognitoidentityprovider.ListUserPoolClientsInput{
		UserPoolId: aws.String(poolId),
	}
	listClientsOutput, err := cognito.ListUserPoolClients(listClientsInput)
	require.NoError(t, err, "failed to list user pool clients")

	var clientId string
	for _, client := range listClientsOutput.UserPoolClients {
		if *client.ClientName == clientName {
			clientId = *client.ClientId
			break
		}
	}
	require.NotEmpty(t, clientId,
		fmt.Sprintf("could not find client with name '%s'", clientName))
	return clientId
}
