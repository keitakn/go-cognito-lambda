package main

import (
	"bytes"
	"html/template"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

// SignUpカスタムメッセージのテスト用期待値を作成する
func createExpectedSignUpMessage(sm SignupMessage) (*bytes.Buffer, error) {
	t := template.New("signup-template.html")

	currentDir, _ := os.Getwd()
	templatePath := currentDir + "/signup-template.html"

	templates := template.Must(t.ParseFiles(templatePath))

	var bodyBuffer bytes.Buffer
	err := templates.Execute(&bodyBuffer, sm)
	if err != nil {
		return nil, err
	}

	return &bodyBuffer, nil
}

// ForgotPasswordMessageカスタムメッセージのテスト用期待値を作成する
func createExpectedForgotPasswordMessageMessage(fm ForgotPasswordMessage) (*bytes.Buffer, error) {
	t := template.New("forgot-password-template.html")

	currentDir, _ := os.Getwd()
	templatePath := currentDir + "/forgot-password-template.html"

	templates := template.Must(t.ParseFiles(templatePath))

	var bodyBuffer bytes.Buffer
	err := templates.Execute(&bodyBuffer, fm)
	if err != nil {
		return nil, err
	}

	return &bodyBuffer, nil
}

func TestHandler(t *testing.T) {
	// TriggerSourceが 'CustomMessage_SignUp' の場合はCustomMessageが返却される
	t.Run("Return Signup CustomMessage", func(t *testing.T) {
		cc := &events.CognitoEventUserPoolsCallerContext{
			AWSSDKVersion: "",
			ClientID:      "",
		}

		ch := &events.CognitoEventUserPoolsHeader{
			Version:       "",
			TriggerSource: "CustomMessage_SignUp",
			Region:        "",
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			CallerContext: *cc,
			UserName:      "keitakn",
		}

		ua := map[string]interface{}{
			"sub": "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
		}

		cm := map[string]string{
			"key": "",
		}

		req := &events.CognitoEventUserPoolsCustomMessageRequest{
			UserAttributes:    ua,
			CodeParameter:     "123456789",
			UsernameParameter: "keitakn",
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

		handlerResult, err := handler(*ev)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		sm := SignupMessage{ConfirmUrl: "http://localhost:3900/cognito/signup/confirm?code=123456789&sub=keitakn"}

		body, err := createExpectedSignUpMessage(sm)
		if err != nil {
			t.Fatal("Error Failed to parse HTML Template", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body.String(),
			EmailSubject: "サインアップ メールアドレスの確認をお願いします。",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_ResendCode' の場合はCustomMessageが返却される
	t.Run("Return ResendCode CustomMessage", func(t *testing.T) {
		cc := &events.CognitoEventUserPoolsCallerContext{
			AWSSDKVersion: "",
			ClientID:      "",
		}

		ch := &events.CognitoEventUserPoolsHeader{
			Version:       "",
			TriggerSource: "CustomMessage_ResendCode",
			Region:        "",
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			CallerContext: *cc,
			UserName:      "keitakn",
		}

		ua := map[string]interface{}{
			"sub": "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
		}

		cm := map[string]string{
			"key": "",
		}

		req := &events.CognitoEventUserPoolsCustomMessageRequest{
			UserAttributes:    ua,
			CodeParameter:     "123456789",
			UsernameParameter: "keitakn",
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

		handlerResult, err := handler(*ev)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		sm := SignupMessage{ConfirmUrl: "http://localhost:3900/cognito/signup/confirm?code=123456789&sub=keitakn"}

		body, err := createExpectedSignUpMessage(sm)
		if err != nil {
			t.Fatal("Error Failed to parse HTML Template", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body.String(),
			EmailSubject: "サインアップ メールアドレスの確認をお願いします。",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_ForgotPassword' の場合はCustomMessageが返却される
	t.Run("Return ForgotPassword CustomMessage", func(t *testing.T) {
		cc := &events.CognitoEventUserPoolsCallerContext{
			AWSSDKVersion: "",
			ClientID:      "",
		}

		ch := &events.CognitoEventUserPoolsHeader{
			Version:       "",
			TriggerSource: "CustomMessage_ForgotPassword",
			Region:        "",
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			CallerContext: *cc,
			UserName:      "keitakn",
		}

		ua := map[string]interface{}{
			"sub": "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
		}

		cm := map[string]string{
			"key": "",
		}

		req := &events.CognitoEventUserPoolsCustomMessageRequest{
			UserAttributes:    ua,
			CodeParameter:     "123456789",
			UsernameParameter: "keitakn",
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

		handlerResult, err := handler(*ev)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		fm := ForgotPasswordMessage{ConfirmUrl: "http://localhost:3900/cognito/password/reset/confirm?code=123456789&sub=keitakn"}

		body, err := createExpectedForgotPasswordMessageMessage(fm)
		if err != nil {
			t.Fatal("Error Failed to parse HTML Template", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body.String(),
			EmailSubject: "パスワードをリセットします。",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_SignUp' だがUserPoolIDが一致しないのでDefaultのメッセージが返却される
	t.Run("Return Signup DefaultMessage Because the UserPoolID doesn't match", func(t *testing.T) {
		cc := &events.CognitoEventUserPoolsCallerContext{
			AWSSDKVersion: "",
			ClientID:      "",
		}

		ch := &events.CognitoEventUserPoolsHeader{
			Version:       "",
			TriggerSource: "CustomMessage_SignUp",
			Region:        "",
			UserPoolID:    "OtherUserPoolID",
			CallerContext: *cc,
			UserName:      "keitakn",
		}

		ua := map[string]interface{}{
			"sub": "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
		}

		cm := map[string]string{
			"key": "",
		}

		req := &events.CognitoEventUserPoolsCustomMessageRequest{
			UserAttributes:    ua,
			CodeParameter:     "123456789",
			UsernameParameter: "keitakn",
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

		handlerResult, err := handler(*ev)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "SMSMessage",
			EmailMessage: "EmailMessage",
			EmailSubject: "EmailSubject",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_ResendCode' だがUserPoolIDが一致しないのでDefaultのメッセージが返却される
	t.Run("Return ResendCode DefaultMessage Because the UserPoolID doesn't match", func(t *testing.T) {
		cc := &events.CognitoEventUserPoolsCallerContext{
			AWSSDKVersion: "",
			ClientID:      "",
		}

		ch := &events.CognitoEventUserPoolsHeader{
			Version:       "",
			TriggerSource: "CustomMessage_ResendCode",
			Region:        "",
			UserPoolID:    "OtherUserPoolID",
			CallerContext: *cc,
			UserName:      "keitakn",
		}

		ua := map[string]interface{}{
			"sub": "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
		}

		cm := map[string]string{
			"key": "",
		}

		req := &events.CognitoEventUserPoolsCustomMessageRequest{
			UserAttributes:    ua,
			CodeParameter:     "123456789",
			UsernameParameter: "keitakn",
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

		handlerResult, err := handler(*ev)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "SMSMessage",
			EmailMessage: "EmailMessage",
			EmailSubject: "EmailSubject",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_ForgotPassword' だがUserPoolIDが一致しないのでDefaultのメッセージが返却される
	t.Run("Return ForgotPassword DefaultMessage Because the UserPoolID doesn't match", func(t *testing.T) {
		cc := &events.CognitoEventUserPoolsCallerContext{
			AWSSDKVersion: "",
			ClientID:      "",
		}

		ch := &events.CognitoEventUserPoolsHeader{
			Version:       "",
			TriggerSource: "CustomMessage_ForgotPassword",
			Region:        "",
			UserPoolID:    "OtherUserPoolID",
			CallerContext: *cc,
			UserName:      "keitakn",
		}

		ua := map[string]interface{}{
			"sub": "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
		}

		cm := map[string]string{
			"key": "",
		}

		req := &events.CognitoEventUserPoolsCustomMessageRequest{
			UserAttributes:    ua,
			CodeParameter:     "123456789",
			UsernameParameter: "keitakn",
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

		handlerResult, err := handler(*ev)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "SMSMessage",
			EmailMessage: "EmailMessage",
			EmailSubject: "EmailSubject",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが指定した値以外の場合はDefaultのメッセージが返却される
	t.Run("Return DefaultMessage Because the TriggerSource is not a specified value", func(t *testing.T) {
		cc := &events.CognitoEventUserPoolsCallerContext{
			AWSSDKVersion: "",
			ClientID:      "",
		}

		ch := &events.CognitoEventUserPoolsHeader{
			Version:       "",
			TriggerSource: "Unknown",
			Region:        "",
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			CallerContext: *cc,
			UserName:      "keitakn",
		}

		ua := map[string]interface{}{
			"sub": "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
		}

		cm := map[string]string{
			"key": "",
		}

		req := &events.CognitoEventUserPoolsCustomMessageRequest{
			UserAttributes:    ua,
			CodeParameter:     "123456789",
			UsernameParameter: "keitakn",
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

		handlerResult, err := handler(*ev)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "SMSMessage",
			EmailMessage: "EmailMessage",
			EmailSubject: "EmailSubject",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})
}
