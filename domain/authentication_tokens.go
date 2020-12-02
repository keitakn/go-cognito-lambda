package domain

type AuthenticationTokens struct {
	Token          string `dynamo:"Token"`
	CognitoSub     string `dynamo:"CognitoSub"`
	SubscribeNews  bool   `dynamo:"SubscribeNews"`
	ExpirationTime int64  `dynamo:"ExpirationTime"`
}
