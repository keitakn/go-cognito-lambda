package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

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

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: "メールアドレスを検証するには、次のリンクをクリックしてください。 http://localhost:3900/cognito/signup/confirm?code=123456789&sub=keitakn",
			EmailSubject: "サインアップ メールアドレスの確認をお願いします。",
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

	// TriggerSourceが指定した値以外の場合はDefaultのメッセージが返却される
	t.Run("Return Signup DefaultMessage Because the TriggerSource is not a specified value", func(t *testing.T) {
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
