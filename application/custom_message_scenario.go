package application

import (
	"bytes"
	"html/template"

	"github.com/keitakn/go-cognito-lambda/domain"
)

type CustomMessageScenario struct {
	Templates                     *template.Template
	AuthenticationTokenRepository domain.AuthenticationTokenRepository
	AuthenticationTokensCreator   domain.AuthenticationTokensCreator
}

type SignUpMessageBuildParams struct {
	Code          string
	SubscribeNews bool
}

type ForgotPasswordMessageBuildParams struct {
	Code string
}

type BuildMessage struct {
	ConfirmUrl string
}

func (s *CustomMessageScenario) BuildSignupMessage(p SignUpMessageBuildParams) (body string, string error) {
	authenticationTokens, err := s.AuthenticationTokensCreator.Create()
	if err != nil {
		return "", err
	}

	if err := s.AuthenticationTokenRepository.Create(*authenticationTokens); err != nil {
		return "", err
	}

	m := BuildMessage{
		ConfirmUrl: "http://localhost:3900/cognito/signup/confirm?code=" + p.Code + "&sub=" + s.AuthenticationTokensCreator.CognitoSub + "&authenticationToken=" + authenticationTokens.Token,
	}

	var bodyBuffer bytes.Buffer
	if err := s.Templates.ExecuteTemplate(&bodyBuffer, "signup-template.html", m); err != nil {
		return "", err
	}

	return bodyBuffer.String(), nil
}

func (s *CustomMessageScenario) BuildForgotPasswordMessage(p ForgotPasswordMessageBuildParams) (body string, string error) {
	m := BuildMessage{
		ConfirmUrl: "http://localhost:3900/cognito/password/reset/confirm?code=" + p.Code + "&sub=" + s.AuthenticationTokensCreator.CognitoSub,
	}

	var bodyBuffer bytes.Buffer
	if err := s.Templates.ExecuteTemplate(&bodyBuffer, "forgot-password-template.html", m); err != nil {
		return "", err
	}

	return bodyBuffer.String(), nil
}
