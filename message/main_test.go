package main

import (
	"github.com/keitakn/go-cognito-lambda/infrastructure"
	"github.com/keitakn/go-cognito-lambda/test"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestMain(m *testing.M) {
	dynamodbClientCreator := infrastructure.DynamodbClientCreator{}

	dynamodb := dynamodbClientCreator.CreateTestClient()

	dynamodbHelper := test.DynamodbHelper{Dynamodb: dynamodb}

	if err := dynamodbHelper.CreateTestAuthenticationTokensTable(); err != nil {
		log.Fatal(err)
	}

	status := m.Run()

	if err := dynamodbHelper.DeleteTestAuthenticationTokensTable(); err != nil {
		log.Fatal(err)
	}

	os.Exit(status)
}

func TestHandler(t *testing.T) {
	// TriggerSourceが 'CustomMessage_SignUp' の場合はCustomMessageが返却される
	t.Run("Return Signup CustomMessage", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_SignUp",
			UserPoolId:    os.Getenv("TARGET_USER_POOL_ID"),
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := "サインアップ メールアドレスの確認をお願いします。"

		if reflect.DeepEqual(handlerResult.Response.EmailSubject, expected) == false {
			t.Error("\nActually: ", handlerResult.Response.EmailSubject, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_ResendCode' の場合はCustomMessageが返却される
	t.Run("Return ResendCode CustomMessage", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_ResendCode",
			UserPoolId:    os.Getenv("TARGET_USER_POOL_ID"),
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := "サインアップ メールアドレスの確認をお願いします。"

		if reflect.DeepEqual(handlerResult.Response.EmailSubject, expected) == false {
			t.Error("\nActually: ", handlerResult.Response.EmailSubject, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_ForgotPassword' の場合はCustomMessageが返却される
	t.Run("Return ForgotPassword CustomMessage", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_ForgotPassword",
			UserPoolId:    os.Getenv("TARGET_USER_POOL_ID"),
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := "パスワードをリセットします。"

		if reflect.DeepEqual(handlerResult.Response.EmailSubject, expected) == false {
			t.Error("\nActually: ", handlerResult.Response.EmailSubject, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_SignUp' だがUserPoolIDが一致しないのでDefaultのメッセージが返却される
	t.Run("Return Signup DefaultMessage Because the UserPoolId doesn't match", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_SignUp",
			UserPoolId:    "OtherUserPoolID",
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
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
	t.Run("Return ResendCode DefaultMessage Because the UserPoolId doesn't match", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_ResendCode",
			UserPoolId:    "OtherUserPoolID",
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
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
	t.Run("Return ForgotPassword DefaultMessage Because the UserPoolId doesn't match", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_ForgotPassword",
			UserPoolId:    "OtherUserPoolID",
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
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
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "Unknown",
			UserPoolId:    os.Getenv("TARGET_USER_POOL_ID"),
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
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
