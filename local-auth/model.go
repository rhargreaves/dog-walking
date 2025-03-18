package main

type AuthorizerResponse struct {
	PrincipalID    string `json:"principalId"`
	PolicyDocument struct {
		Version   string      `json:"Version"`
		Statement []Statement `json:"Statement"`
	} `json:"policyDocument"`
	Context map[string]interface{} `json:"context"`
}

type Statement struct {
	Action   string `json:"Action"`
	Effect   string `json:"Effect"`
	Resource string `json:"Resource"`
}
