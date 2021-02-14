package domain

type CognitoIdToken struct {
	Sub           string `json:"sub"`
	Aud           string `json:"aud"`
	Iss           string `json:"iss"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
}
