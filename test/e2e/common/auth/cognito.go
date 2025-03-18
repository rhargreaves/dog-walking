package auth

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func GetCognitoJWT() string {
	poolName := os.Getenv("COGNITO_USER_POOL_NAME")
	if poolName == "" {
		log.Fatal("COGNITO_USER_POOL_NAME environment variable is required")
	}
	clientName := os.Getenv("COGNITO_CLIENT_NAME")
	if clientName == "" {
		log.Fatal("COGNITO_CLIENT_NAME environment variable is required")
	}

	sess := session.Must(session.NewSession())
	cognito := cognitoidentityprovider.New(sess)
	userPoolId := findUserPoolByName(cognito, poolName)
	clientId := findClientByName(cognito, userPoolId, clientName)

	username := os.Getenv("COGNITO_USERNAME")
	if username == "" {
		log.Fatal("COGNITO_USERNAME environment variable is required")
	}
	password := os.Getenv("COGNITO_PASSWORD")
	if password == "" {
		log.Fatal("COGNITO_PASSWORD environment variable is required")
	}
	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String(clientId),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
	}

	authOutput, err := cognito.InitiateAuth(authInput)
	if err != nil {
		log.Fatal("failed to authenticate with Cognito:", err)
	}
	if authOutput.AuthenticationResult == nil {
		log.Fatal("no authentication result received")
	}
	return *authOutput.AuthenticationResult.IdToken
}

func findUserPoolByName(cognito *cognitoidentityprovider.CognitoIdentityProvider, poolName string) string {
	listPoolsInput := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int64(1),
	}
	listPoolsOutput, err := cognito.ListUserPools(listPoolsInput)
	if err != nil {
		log.Fatal("failed to list user pools:", err)
	}

	var userPoolId string
	for _, pool := range listPoolsOutput.UserPools {
		if *pool.Name == poolName {
			userPoolId = *pool.Id
			break
		}
	}
	if userPoolId == "" {
		log.Fatalf("could not find user pool with name '%s'", poolName)
	}
	return userPoolId
}

func findClientByName(cognito *cognitoidentityprovider.CognitoIdentityProvider, poolId string, clientName string) string {
	listClientsInput := &cognitoidentityprovider.ListUserPoolClientsInput{
		UserPoolId: aws.String(poolId),
	}
	listClientsOutput, err := cognito.ListUserPoolClients(listClientsInput)
	if err != nil {
		log.Fatal("failed to list user pool clients:", err)
	}

	var clientId string
	for _, client := range listClientsOutput.UserPoolClients {
		if *client.ClientName == clientName {
			clientId = *client.ClientId
			break
		}
	}
	if clientId == "" {
		log.Fatalf("could not find client with name '%s'", clientName)
	}
	return clientId
}
