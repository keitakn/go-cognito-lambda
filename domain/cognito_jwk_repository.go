package domain

import "github.com/lestrrat-go/jwx/jwk"

type CognitoJwkRepository interface {
	Fetch() (jwk.Set, error)
}
