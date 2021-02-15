package repository

import (
	"github.com/keitakn/go-cognito-lambda/domain"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// https://github.com/lestrrat-go/jwx を利用して実装
type JwxCognitoJwtTokenRepository struct {
	Iss    string
	Aud    string
	JwkSet jwk.Set
}

func (r *JwxCognitoJwtTokenRepository) ParseAndValidateIdToken(tokenStr string) (*domain.CognitoIdToken, error) {
	tokenObj, err := jwt.ParseString(
		tokenStr,
		jwt.WithKeySet(r.JwkSet),
	)

	// TokenのParseに失敗した場合はエラーを返す
	if err != nil {
		return nil, err
	}

	err = jwt.Validate(
		tokenObj,
		jwt.WithIssuer(r.Iss),
		jwt.WithAudience(r.Aud),
		jwt.WithClaimValue("token_use", "id"),
	)

	// Tokenの検証に失敗した場合はエラーを返す
	if err != nil {
		return nil, err
	}

	idToken := &domain.CognitoIdToken{
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

func (
	r *JwxCognitoJwtTokenRepository,
) ParseAndValidateAccessToken(
	tokenStr string,
) (
	*domain.CognitoAccessToken, error,
) {
	tokenObj, err := jwt.ParseString(
		tokenStr,
		jwt.WithKeySet(r.JwkSet),
	)

	// TokenのParseに失敗した場合はエラーを返す
	if err != nil {
		return nil, err
	}

	err = jwt.Validate(
		tokenObj,
		jwt.WithIssuer(r.Iss),
		jwt.WithClaimValue("token_use", "access"),
		jwt.WithClaimValue("client_id", r.Aud),
	)

	// Tokenの検証に失敗した場合はエラーを返す
	if err != nil {
		return nil, err
	}

	accessToken := &domain.CognitoAccessToken{
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
