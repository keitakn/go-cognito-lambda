package main

import (
	"bytes"
	"github.com/aws/aws-lambda-go/events"
	"html/template"
	"os"
)

// テスト用の期待値を作成する
func createExpectedSignUpMessage(m SignUpMessage) (*bytes.Buffer, error) {
	t := template.New("signup-template.html")

	currentDir, _ := os.Getwd()
	templatePath := currentDir + "/signup-template.html"

	templates := template.Must(t.ParseFiles(templatePath))

	var bodyBuffer bytes.Buffer
	err := templates.Execute(&bodyBuffer, m)
	if err != nil {
		return nil, err
	}

	return &bodyBuffer, nil
}

// ForgotPasswordMessageカスタムメッセージのテスト用期待値を作成する
func createExpectedForgotPasswordMessageMessage(m ForgotPasswordMessage) (*bytes.Buffer, error) {
	t := template.New("forgot-password-template.html")

	currentDir, _ := os.Getwd()
	templatePath := currentDir + "/forgot-password-template.html"

	templates := template.Must(t.ParseFiles(templatePath))

	var bodyBuffer bytes.Buffer
	err := templates.Execute(&bodyBuffer, m)
	if err != nil {
		return nil, err
	}

	return &bodyBuffer, nil
}

type createUserPoolsCustomMessageEventParams struct {
	TriggerSource string
	UserPoolId    string
	UserName      string
	Sub           string
	CodeParameter string
}

func createUserPoolsCustomMessageEvent(p *createUserPoolsCustomMessageEventParams) *events.CognitoEventUserPoolsCustomMessage {
	cc := &events.CognitoEventUserPoolsCallerContext{
		AWSSDKVersion: "",
		ClientID:      "",
	}

	ch := &events.CognitoEventUserPoolsHeader{
		Version:       "",
		TriggerSource: p.TriggerSource,
		Region:        "",
		UserPoolID:    p.UserPoolId,
		CallerContext: *cc,
		UserName:      p.UserName,
	}

	ua := map[string]interface{}{
		"sub": p.Sub,
	}

	cm := map[string]string{
		"key": "",
	}

	req := &events.CognitoEventUserPoolsCustomMessageRequest{
		UserAttributes:    ua,
		CodeParameter:     p.CodeParameter,
		UsernameParameter: p.UserName,
		ClientMetadata:    cm,
	}

	res := &events.CognitoEventUserPoolsCustomMessageResponse{
		SMSMessage:   "SMSMessage",
		EmailMessage: "EmailMessage",
		EmailSubject: "EmailSubject",
	}

	ev := &events.CognitoEventUserPoolsCustomMessage{
		CognitoEventUserPoolsHeader: *ch,
		Request:                     *req,
		Response:                    *res,
	}

	return ev
}
