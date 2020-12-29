package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// 「認証チャレンジの作成 Lambda」
// event.Response.PublicChallengeParameters["challenge"] の値はクライアントに渡される
// event.Response.PrivateChallengeParameters["answer"] の値は「認証チャレンジレスポンスの確認 Lambda」（authchallenge/verify/main.go）に渡される
// https://docs.aws.amazon.com/ja_jp/cognito/latest/developerguide/user-pool-lambda-create-auth-challenge.html
func Handler(
	event events.CognitoEventUserPoolsCreateAuthChallenge,
) (events.CognitoEventUserPoolsCreateAuthChallenge, error) {
	event.Response.PublicChallengeParameters = map[string]string{}
	event.Response.PrivateChallengeParameters = map[string]string{}

	if event.Request.ChallengeName == "CUSTOM_CHALLENGE" {
		event.Response.PublicChallengeParameters["challenge"] = "AuthenticationTokens"
		event.Response.PrivateChallengeParameters["answer"] = event.UserName
		event.Response.ChallengeMetadata = "CHALLENGE_AND_RESPONSE" + string(rune(len(event.Request.Session)))
	}

	return event, nil
}

func main() {
	lambda.Start(Handler)
}
