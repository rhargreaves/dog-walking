package main

import "github.com/aws/aws-lambda-go/events"

func response(principalID string, methodArn string, context map[string]any,
	effect string) events.APIGatewayV2CustomAuthorizerIAMPolicyResponse {
	return events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{
		PrincipalID: principalID,
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{methodArn},
				},
			},
		},
		Context: context,
	}
}

func errorResponse(errorMessage string,
	methodArn string) events.APIGatewayV2CustomAuthorizerIAMPolicyResponse {
	return response("", methodArn, map[string]any{
		"error": errorMessage,
	}, "Deny")
}

func authorisedResponse(principalID string, methodArn string, email string,
	groups []string) events.APIGatewayV2CustomAuthorizerIAMPolicyResponse {
	return response(principalID, methodArn, map[string]any{
		"userId": principalID,
		"email":  email,
		"groups": groups,
	}, "Allow")
}
