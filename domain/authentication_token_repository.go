package domain

type AuthenticationTokenRepository interface {
	Create(item AuthenticationTokens) error
	FindByToken(token string) (*AuthenticationTokens, error)
}
