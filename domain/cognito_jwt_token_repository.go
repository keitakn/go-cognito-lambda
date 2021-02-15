package domain

type CognitoJwtTokenRepository interface {
	ParseAndValidateIdToken(tokenStr string) (*CognitoIdToken, error)
	ParseAndValidateAccessToken(tokenStr string) (*CognitoAccessToken, error)
}
