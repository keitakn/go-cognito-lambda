package domain

import (
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type CognitoJwtToken struct {
	Iss    string
	Aud    string
	JwkSet jwk.Set
}

func (c *CognitoJwtToken) ParseAndValidateIdToken(tokenStr string) (*CognitoIdToken, error) {
	tokenObj, err := jwt.ParseString(
		tokenStr,
		jwt.WithKeySet(c.JwkSet),
	)

	// TokenのParseに失敗した場合はエラーを返す
	if err != nil {
		return nil, err
	}

	err = jwt.Validate(
		tokenObj,
		jwt.WithIssuer(c.Iss),
		jwt.WithAudience(c.Aud),
		jwt.WithClaimValue("token_use", "id"),
	)

	// Tokenの検証に失敗した場合はエラーを返す
	if err != nil {
		return nil, err
	}

	idToken := &CognitoIdToken{
		Sub: tokenObj.Subject(),
		Aud: tokenObj.Audience()[0],
		Iss: tokenObj.Issuer(),
	}

	if val, ok := tokenObj.PrivateClaims()["email"]; ok {
		idToken.Email = val.(string)
	}

	if val, ok := tokenObj.PrivateClaims()["email_verified"]; ok {
		idToken.EmailVerified = val.(bool)
	}

	return idToken, nil
}

func (c *CognitoJwtToken) ParseAndValidateAccessToken(tokenStr string) (*CognitoAccessToken, error) {
	tokenObj, err := jwt.ParseString(
		tokenStr,
		jwt.WithKeySet(c.JwkSet),
	)

	// TokenのParseに失敗した場合はエラーを返す
	if err != nil {
		return nil, err
	}

	err = jwt.Validate(
		tokenObj,
		jwt.WithIssuer(c.Iss),
		jwt.WithClaimValue("token_use", "access"),
		jwt.WithClaimValue("client_id", c.Aud),
	)

	// Tokenの検証に失敗した場合はエラーを返す
	if err != nil {
		return nil, err
	}

	accessToken := &CognitoAccessToken{
		Sub: tokenObj.Subject(),
		Iss: tokenObj.Issuer(),
	}

	if val, ok := tokenObj.PrivateClaims()["scope"]; ok {
		accessToken.Scope = val.(string)
	}

	if val, ok := tokenObj.PrivateClaims()["client_id"]; ok {
		accessToken.ClientId = val.(string)
	}

	return accessToken, nil
}
