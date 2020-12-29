package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// 「認証チャレンジの定義 Lambda」 この後に 「認証チャレンジの作成 Lambda」（authchallenge/create/main.go）が呼ばれる
// https://docs.aws.amazon.com/ja_jp/cognito/latest/developerguide/user-pool-lambda-define-auth-challenge.html
func Handler(
	event events.CognitoEventUserPoolsDefineAuthChallenge,
) (events.CognitoEventUserPoolsDefineAuthChallenge, error) {
	targetUserPoolId := os.Getenv("TARGET_USER_POOL_ID")
	if targetUserPoolId != event.UserPoolID {
		return event, nil
	}

	sessionCnt := len(event.Request.Session)

	if sessionCnt == 0 {
		// パスワードを省略して、signInメソッドを呼ぶとここが呼ばれ、クライアント側で Custom authentication flow に入る
		// https://docs.amplify.aws/lib/auth/switch-auth/q/platform/js#custom_auth-flow
		event.Response.ChallengeName = "CUSTOM_CHALLENGE"
		event.Response.FailAuthentication = false
		event.Response.IssueTokens = false
	} else if sessionCnt > 0 && event.Request.Session[sessionCnt-1].ChallengeResult {
		// カスタム認証に成功した場合
		event.Response.FailAuthentication = false
		event.Response.IssueTokens = true
	} else {
		// カスタム認証に失敗した場合
		event.Response.FailAuthentication = true
		event.Response.IssueTokens = false
	}

	return event, nil
}

func main() {
	lambda.Start(Handler)
}
