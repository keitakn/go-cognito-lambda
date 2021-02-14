package domain

type CognitoAccessToken struct {
	Sub      string `json:"sub"`
	Iss      string `json:"iss"`
	Scope    string `json:"scope"`
	ClientId string `json:"clientId"`
}
