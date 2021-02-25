package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestMain(m *testing.M) {
	status := m.Run()

	os.Exit(status)
}

//nolint:funlen
func TestHandler(t *testing.T) {
	const expectedEmailVerified = "true"

	const expectedFinalUserStatus = "CONFIRMED"

	const expectedMessageAction = "SUPPRESS"

	//nolint:dupl
	t.Run("Successful user migration TriggerSource is UserMigration_Authentication", func(t *testing.T) {
		eventHeader := &events.CognitoEventUserPoolsHeader{
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			TriggerSource: "UserMigration_Authentication",
			UserName:      "keita.koga.work+migrateuser@gmail.com",
		}

		event := &events.CognitoEventUserPoolsMigrateUser{
			CognitoEventUserPoolsHeader: *eventHeader,
		}

		handlerResult, err := Handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expectedMigrateUser := make(map[string]string)
		expectedMigrateUser["email"] = event.UserName
		expectedMigrateUser["email_verified"] = expectedEmailVerified

		if reflect.DeepEqual(
			handlerResult.CognitoEventUserPoolsMigrateUserResponse.UserAttributes, expectedMigrateUser,
		) == false {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.UserAttributes,
				"\nExpected: ",
				expectedMigrateUser,
			)
		}

		if handlerResult.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus != expectedFinalUserStatus {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus,
				"\nExpected: ",
				expectedFinalUserStatus,
			)
		}

		if handlerResult.CognitoEventUserPoolsMigrateUserResponse.MessageAction != expectedMessageAction {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.MessageAction,
				"\nExpected: ",
				expectedMessageAction,
			)
		}
	})

	//nolint:dupl
	t.Run("Successful user migration TriggerSource is UserMigration_ForgotPassword", func(t *testing.T) {
		eventHeader := &events.CognitoEventUserPoolsHeader{
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			TriggerSource: "UserMigration_ForgotPassword",
			UserName:      "keita.koga.work+migrateuser@gmail.com",
		}

		event := &events.CognitoEventUserPoolsMigrateUser{
			CognitoEventUserPoolsHeader: *eventHeader,
		}

		handlerResult, err := Handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expectedMigrateUser := make(map[string]string)
		expectedMigrateUser["email"] = event.UserName
		expectedMigrateUser["email_verified"] = expectedEmailVerified

		if reflect.DeepEqual(
			handlerResult.CognitoEventUserPoolsMigrateUserResponse.UserAttributes, expectedMigrateUser,
		) == false {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.UserAttributes,
				"\nExpected: ",
				expectedMigrateUser,
			)
		}

		if handlerResult.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus != expectedFinalUserStatus {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus,
				"\nExpected: ",
				expectedFinalUserStatus,
			)
		}

		if handlerResult.CognitoEventUserPoolsMigrateUserResponse.MessageAction != expectedMessageAction {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.MessageAction,
				"\nExpected: ",
				expectedMessageAction,
			)
		}
	})

	t.Run("The user will not be migrated, Because the UserPoolID is different", func(t *testing.T) {
		eventHeader := &events.CognitoEventUserPoolsHeader{
			UserPoolID:    "UNKNOWN_USER_POOL_ID",
			TriggerSource: "UserMigration_ForgotPassword",
			UserName:      "keita.koga.work+migrateuser@gmail.com",
		}

		event := &events.CognitoEventUserPoolsMigrateUser{
			CognitoEventUserPoolsHeader: *eventHeader,
		}

		handlerResult, err := Handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expectedFinalUserStatus := ""
		if handlerResult.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus != expectedFinalUserStatus {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus,
				"\nExpected: ",
				expectedFinalUserStatus,
			)
		}

		expectedMessageAction := ""
		if handlerResult.CognitoEventUserPoolsMigrateUserResponse.MessageAction != expectedMessageAction {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.MessageAction,
				"\nExpected: ",
				expectedMessageAction,
			)
		}
	})

	t.Run("The user will not be migrated, Because the TriggerSource is different", func(t *testing.T) {
		eventHeader := &events.CognitoEventUserPoolsHeader{
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			TriggerSource: "Unknown_TriggerSource",
			UserName:      "keita.koga.work+migrateuser@gmail.com",
		}

		event := &events.CognitoEventUserPoolsMigrateUser{
			CognitoEventUserPoolsHeader: *eventHeader,
		}

		handlerResult, err := Handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expectedFinalUserStatus := ""
		if handlerResult.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus != expectedFinalUserStatus {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus,
				"\nExpected: ",
				expectedFinalUserStatus,
			)
		}

		expectedMessageAction := ""
		if handlerResult.CognitoEventUserPoolsMigrateUserResponse.MessageAction != expectedMessageAction {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.MessageAction,
				"\nExpected: ",
				expectedMessageAction,
			)
		}
	})

	t.Run("The user will not be migrated, Because the UserName is different", func(t *testing.T) {
		eventHeader := &events.CognitoEventUserPoolsHeader{
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			TriggerSource: "UserMigration_Authentication",
			UserName:      "keita.koga.work@gmail.com",
		}

		event := &events.CognitoEventUserPoolsMigrateUser{
			CognitoEventUserPoolsHeader: *eventHeader,
		}

		handlerResult, err := Handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expectedFinalUserStatus := ""
		if handlerResult.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus != expectedFinalUserStatus {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus,
				"\nExpected: ",
				expectedFinalUserStatus,
			)
		}

		expectedMessageAction := ""
		if handlerResult.CognitoEventUserPoolsMigrateUserResponse.MessageAction != expectedMessageAction {
			t.Error(
				"\nActually: ",
				handlerResult.CognitoEventUserPoolsMigrateUserResponse.MessageAction,
				"\nExpected: ",
				expectedMessageAction,
			)
		}
	})
}
