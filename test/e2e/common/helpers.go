package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func BaseUrl() string {
	return os.Getenv("API_BASE_URL")
}

func NewAuthedRequest(t *testing.T, method, endpoint string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, BaseUrl()+endpoint, body)
	require.NoError(t, err, "Failed to create "+method+" request")
	req.Header.Set("Authorization", "Bearer "+createTestJWT(t))
	return req
}

func Get(t *testing.T, endpoint string) *http.Response {
	req := NewAuthedRequest(t, http.MethodGet, endpoint, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to fetch URL")
	return resp
}

func PostBytes(t *testing.T, endpoint string, body []byte) *http.Response {
	req := NewAuthedRequest(t, http.MethodPost, endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	require.NoError(t, err, "Failed to POST to URL")
	return resp
}

func PostJson(t *testing.T, endpoint string, body interface{}) *http.Response {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err, "Failed to marshal body")
	return PostBytes(t, endpoint, jsonBody)
}

func PutBytes(t *testing.T, endpoint string, body []byte) *http.Response {
	req := NewAuthedRequest(t, http.MethodPut, endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to PUT to URL")
	return resp
}

func PutJson(t *testing.T, endpoint string, body interface{}) *http.Response {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err, "Failed to marshal body")
	return PutBytes(t, endpoint, jsonBody)
}

func Delete(t *testing.T, endpoint string) *http.Response {
	req := NewAuthedRequest(t, http.MethodDelete, endpoint, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to DELETE to URL")
	return resp
}

func RequireStatus(t *testing.T, resp *http.Response, expectedStatus int) {
	assert.Equal(t, expectedStatus, resp.StatusCode,
		"Expected status code %d", expectedStatus)

	if resp.StatusCode != expectedStatus {
		logBody(t, resp)
	}
}

func DecodeJSON(t *testing.T, resp *http.Response, target interface{}) {
	err := json.NewDecoder(resp.Body).Decode(target)
	require.NoError(t, err, "Failed to decode response body")
}

func logBody(t *testing.T, resp *http.Response) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	t.Errorf("Body: %s", string(bodyBytes))
}

func SkipIfLocal(t *testing.T) {
	if strings.HasPrefix(BaseUrl(), "http://sam:") {
		t.Skip("Skipping test on local environment")
	}
}

func findUserPoolByName(t *testing.T, cognito *cognitoidentityprovider.CognitoIdentityProvider, poolName string) string {
	listPoolsInput := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int64(1),
	}
	listPoolsOutput, err := cognito.ListUserPools(listPoolsInput)
	require.NoError(t, err, "Failed to list user pools")

	var userPoolId string
	for _, pool := range listPoolsOutput.UserPools {
		if *pool.Name == poolName {
			userPoolId = *pool.Id
			break
		}
	}
	require.NotEmpty(t, userPoolId,
		fmt.Sprintf("Could not find user pool with name '%s'", poolName))
	return userPoolId
}

func findClientByName(t *testing.T, cognito *cognitoidentityprovider.CognitoIdentityProvider, poolId string, clientName string) string {
	listClientsInput := &cognitoidentityprovider.ListUserPoolClientsInput{
		UserPoolId: aws.String(poolId),
	}
	listClientsOutput, err := cognito.ListUserPoolClients(listClientsInput)
	require.NoError(t, err, "Failed to list user pool clients")

	var clientId string
	for _, client := range listClientsOutput.UserPoolClients {
		if *client.ClientName == clientName {
			clientId = *client.ClientId
			break
		}
	}
	require.NotEmpty(t, clientId,
		fmt.Sprintf("Could not find client with name '%s'", clientName))
	return clientId
}

func getCognitoToken(t *testing.T) string {
	poolName := os.Getenv("COGNITO_USER_POOL_NAME")
	require.NotEmpty(t, poolName, "COGNITO_USER_POOL_NAME environment variable is required")

	clientName := os.Getenv("COGNITO_CLIENT_NAME")
	require.NotEmpty(t, clientName, "COGNITO_CLIENT_NAME environment variable is required")

	sess := session.Must(session.NewSession())
	cognito := cognitoidentityprovider.New(sess)
	userPoolId := findUserPoolByName(t, cognito, poolName)
	clientId := findClientByName(t, cognito, userPoolId, clientName)

	// Get credentials from environment
	username := os.Getenv("COGNITO_USERNAME")
	password := os.Getenv("COGNITO_PASSWORD")
	require.NotEmpty(t, username, "COGNITO_USERNAME environment variable is required")
	require.NotEmpty(t, password, "COGNITO_PASSWORD environment variable is required")

	// Initiate auth
	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String(clientId),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
	}

	authOutput, err := cognito.InitiateAuth(authInput)
	require.NoError(t, err, "Failed to authenticate with Cognito")
	require.NotNil(t, authOutput.AuthenticationResult, "No authentication result received")

	return *authOutput.AuthenticationResult.IdToken
}

func createTestJWT(t *testing.T) string {
	if os.Getenv("USE_REAL_COGNITO") == "true" {
		return getCognitoToken(t)
	}

	// Fallback to local test JWT
	claims := jwt.MapClaims{
		"sub":            "test-user-id",
		"email":          "test@example.com",
		"cognito:groups": []string{"Users"},
		"exp":            time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("LOCAL_JWT_SECRET")))
	require.NoError(t, err, "Failed to create JWT")
	fmt.Println("tokenString:", tokenString)
	return tokenString
}
