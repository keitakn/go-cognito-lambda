package domain

type AuthenticationTokens struct {
	Token          string `dynamodbav:"Token"`
	CognitoSub     string `dynamodbav:"CognitoSub"`
	SubscribeNews  bool   `dynamodbav:"SubscribeNews"`
	ExpirationTime int64  `dynamodbav:"ExpirationTime"`
}
