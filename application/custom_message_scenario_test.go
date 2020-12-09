package application

import (
	"github.com/keitakn/go-cognito-lambda/domain"
	"html/template"
	"os"
	"reflect"
	"testing"
	"time"
)

var templates *template.Template

func TestMain(m *testing.M) {
	signupTemplatePath := "../message/signup-template.html"
	forgotPasswordTemplatePath := "../message/forgot-password-template.html"

	templates = template.Must(template.ParseFiles(signupTemplatePath, forgotPasswordTemplatePath))

	status := m.Run()

	os.Exit(status)
}

func TestHandler(t *testing.T) {
	t.Run("Successful BuildSignupMessage", func(t *testing.T) {
		tokensCreator := domain.AuthenticationTokensCreator{
			Token:         "aaaaaaaa-1111-2222-3333-bbbbbbbbbbbb",
			CognitoSub:    "0ef53af5-4eb9-4d2b-a939-8cb9d795512b",
			SubscribeNews: true,
			Time:          time.Now(),
		}

		tokens, err := tokensCreator.Create()
		if err != nil {
			t.Fatal("Error failed to Generate AuthenticationTokens", err)
		}

		repo := MockAuthenticationTokenRepository{
			Token:          tokens.Token,
			CognitoSub:     tokens.CognitoSub,
			SubscribeNews:  tokens.SubscribeNews,
			ExpirationTime: tokens.ExpirationTime,
		}

		scenario := CustomMessageScenario{
			Templates:                     templates,
			AuthenticationTokenRepository: &repo,
			AuthenticationTokensCreator:   tokensCreator,
		}

		code := "123456"
		res, err := scenario.BuildSignupMessage(SignUpMessageBuildParams{Code: code, SubscribeNews: true})
		if err != nil {
			t.Fatal("Error failed to BuildSignupMessage", err)
		}

		m := BuildMessage{
			ConfirmUrl: "http://localhost:3900/cognito/signup/confirm?code=" + code + "&sub=" + tokensCreator.CognitoSub + "&authenticationToken=" + tokensCreator.Token,
		}

		expected, err := createExpectedSignUpMessage(m)
		if err != nil {
			t.Fatal("Error failed to createExpectedSignUpMessage", err)
		}

		if reflect.DeepEqual(res, expected) == false {
			t.Error("\nActually: ", res, "\nExpected: ", expected)
		}
	})

	t.Run("Successful BuildForgotPasswordMessage", func(t *testing.T) {
		tokensCreator := domain.AuthenticationTokensCreator{
			Token:         "aaaaaaaa-1111-2222-3333-bbbbbbbbbbbb",
			CognitoSub:    "0ef53af5-4eb9-4d2b-a939-8cb9d795512b",
			SubscribeNews: true,
			Time:          time.Now(),
		}

		tokens, err := tokensCreator.Create()
		if err != nil {
			t.Fatal("Error failed to Generate AuthenticationTokens", err)
		}

		repo := MockAuthenticationTokenRepository{
			Token:          tokens.Token,
			CognitoSub:     tokens.CognitoSub,
			SubscribeNews:  tokens.SubscribeNews,
			ExpirationTime: tokens.ExpirationTime,
		}

		scenario := CustomMessageScenario{
			Templates:                     templates,
			AuthenticationTokenRepository: &repo,
			AuthenticationTokensCreator:   tokensCreator,
		}

		code := "987654"
		res, err := scenario.BuildForgotPasswordMessage(ForgotPasswordMessageBuildParams{Code: code})
		if err != nil {
			t.Fatal("Error failed to BuildForgotPasswordMessage", err)
		}

		m := BuildMessage{
			ConfirmUrl: "http://localhost:3900/cognito/password/reset/confirm?code=" + code + "&sub=" + tokensCreator.CognitoSub,
		}

		expected, err := createExpectedForgotPasswordMessageMessage(m)
		if err != nil {
			t.Fatal("Error failed to createExpectedForgotPasswordMessageMessage", err)
		}

		if reflect.DeepEqual(res, expected) == false {
			t.Error("\nActually: ", res, "\nExpected: ", expected)
		}
	})
}
