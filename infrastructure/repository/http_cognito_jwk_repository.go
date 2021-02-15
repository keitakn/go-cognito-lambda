package repository

import (
	"context"
	"fmt"

	"github.com/lestrrat-go/jwx/jwk"
)

type HttpCognitoJwkRepository struct {
	Iss string
}

func (r *HttpCognitoJwkRepository) Fetch() (jwk.Set, error) {
	jwkUrl := fmt.Sprintf("%v/.well-known/jwks.json", r.Iss)
	ctx := context.Background()

	jwkSet, err := jwk.Fetch(ctx, jwkUrl)
	if err != nil {
		return nil, err
	}

	return jwkSet, nil
}
