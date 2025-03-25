const { CognitoIdentityProviderClient, InitiateAuthCommand } = require("@aws-sdk/client-cognito-identity-provider");

function getCognitoToken(userContext, events, done) {
  const client = new CognitoIdentityProviderClient({
    region: userContext.vars.region
  });

  const command = new InitiateAuthCommand({
    AuthFlow: "USER_PASSWORD_AUTH",
    ClientId: userContext.vars.clientId,
    AuthParameters: {
      USERNAME: userContext.vars.username,
      PASSWORD: userContext.vars.password
    }
  });

  client.send(command)
    .then(response => {
      userContext.vars.token = response.AuthenticationResult.IdToken;
      done();
    })
    .catch(error => {
      done(error);
    });
}

module.exports = {
  getCognitoToken
};