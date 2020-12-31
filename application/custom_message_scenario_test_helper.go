package application

import (
	"bytes"
	"html/template"

	"github.com/keitakn/go-cognito-lambda/domain"
)

type MockAuthenticationTokenRepository struct {
	Token          string
	CognitoSub     string
	SubscribeNews  bool
	ExpirationTime int64
}

func (r *MockAuthenticationTokenRepository) Create(item domain.AuthenticationTokens) error {
	return nil
}

func (r *MockAuthenticationTokenRepository) FindByToken(token string) (*domain.AuthenticationTokens, error) {
	return &domain.AuthenticationTokens{
		Token:          token,
		CognitoSub:     r.CognitoSub,
		SubscribeNews:  r.SubscribeNews,
		ExpirationTime: r.ExpirationTime,
	}, nil
}

// SignupMessageカスタムメッセージのテスト用期待値を作成する
func createExpectedSignUpMessage(m BuildMessage) (string, error) {
	t := template.New("signup-template.html")

	templatePath := "../message/signup-template.html"

	templates := template.Must(t.ParseFiles(templatePath))

	var bodyBuffer bytes.Buffer
	err := templates.Execute(&bodyBuffer, m)
	if err != nil {
		return "", err
	}

	return bodyBuffer.String(), nil
}

// ForgotPasswordMessageカスタムメッセージのテスト用期待値を作成する
func createExpectedForgotPasswordMessageMessage(m BuildMessage) (string, error) {
	t := template.New("forgot-password-template.html")

	templatePath := "../message/forgot-password-template.html"

	templates := template.Must(t.ParseFiles(templatePath))

	var bodyBuffer bytes.Buffer
	err := templates.Execute(&bodyBuffer, m)
	if err != nil {
		return "", err
	}

	return bodyBuffer.String(), nil
}
