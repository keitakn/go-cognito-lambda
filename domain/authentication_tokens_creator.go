package domain

import (
	"time"

	"github.com/google/uuid"
)

const expireMinute int64 = 10

type AuthenticationTokensCreator struct {
	Token         string
	CognitoSub    string
	SubscribeNews bool
	Time          time.Time
}

func (c *AuthenticationTokensCreator) Create() (*AuthenticationTokens, error) {
	token := c.Token

	if token == "" {
		randomToken, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		token = randomToken.String()
	}

	expirationTime := c.Time.Add(time.Duration(expireMinute) * time.Minute)

	return &AuthenticationTokens{
		Token:          token,
		CognitoSub:     c.CognitoSub,
		SubscribeNews:  c.SubscribeNews,
		ExpirationTime: expirationTime.Unix(),
	}, nil
}
