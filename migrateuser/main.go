package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// 検証用なのでメールアドレスが特定の値だったら移行対象と見なす
func shouldMigrate(email string) bool {
	return email == "keita.koga.work+migrateuser@gmail.com"
}

func Handler(event events.CognitoEventUserPoolsMigrateUser) (events.CognitoEventUserPoolsMigrateUser, error) {
	if event.UserPoolID != os.Getenv("TARGET_USER_POOL_ID") {
		return event, nil
	}

	// email_verifiedは文字列で設定する必要がある
	emailVerified := "true"

	// 認証が行われた際に呼ばれる
	if event.TriggerSource == "UserMigration_Authentication" {
		// 検証用なのでメールアドレスが特定の値だったら移行対象と見なす
		// 本来はここで移行元の認証システムに認証を行い移行対象かどうかを判断する
		if shouldMigrate(event.UserName) {
			migrateUser := make(map[string]string)
			migrateUser["email"] = event.UserName
			migrateUser["email_verified"] = emailVerified

			event.CognitoEventUserPoolsMigrateUserResponse.UserAttributes = migrateUser
			event.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus = "CONFIRMED"
			event.CognitoEventUserPoolsMigrateUserResponse.MessageAction = "SUPPRESS"

			return event, nil
		}
	}

	// パスワードリセットの認証メール送信時に呼ばれる
	if event.TriggerSource == "UserMigration_ForgotPassword" {
		// 検証用なのでメールアドレスが特定の値だったら移行対象と見なす
		// 本来はここで移行元の認証システムに認証を行い移行対象かどうかを判断する
		if shouldMigrate(event.UserName) {
			migrateUser := make(map[string]string)
			migrateUser["email"] = event.UserName
			migrateUser["email_verified"] = emailVerified

			event.CognitoEventUserPoolsMigrateUserResponse.UserAttributes = migrateUser
			event.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus = "CONFIRMED"
			event.CognitoEventUserPoolsMigrateUserResponse.MessageAction = "SUPPRESS"

			return event, nil
		}
	}

	return event, nil
}

func main() {
	lambda.Start(Handler)
}
