package main

import (
	"html/template"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/keitakn/go-cognito-lambda/application"
	"github.com/keitakn/go-cognito-lambda/domain"
	"github.com/keitakn/go-cognito-lambda/infrastructure"
	"github.com/keitakn/go-cognito-lambda/infrastructure/repository"
)

var templates *template.Template
var db *dynamodb.DynamoDB
var authenticationTokenRepository domain.AuthenticationTokenRepository

//nolint:gochecknoinits
func init() {
	signupTemplatePath := "bin/signup-template.html"
	forgotPasswordTemplatePath := "bin/forgot-password-template.html"

	dynamodbClientCreator := infrastructure.DynamodbClientCreator{}

	if infrastructure.IsTestRun() {
		currentDir, _ := os.Getwd()
		signupTemplatePath = currentDir + "/signup-template.html"
		forgotPasswordTemplatePath = currentDir + "/forgot-password-template.html"

		db = dynamodbClientCreator.CreateTestClient()
	} else {
		db = dynamodbClientCreator.Create()
	}

	templates = template.Must(template.ParseFiles(signupTemplatePath, forgotPasswordTemplatePath))

	authenticationTokenRepository = &repository.DynamodbAuthenticationTokenRepository{Dynamodb: db}
}

func handler(request events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessage, error) {
	targetUserPoolId := os.Getenv("TARGET_USER_POOL_ID")

	// 実行対象のユーザープールからのリクエストではない場合は何もせずにデフォルトのメッセージを返す
	if targetUserPoolId != request.UserPoolID {
		return request, nil
	}

	// サインアップ時に送られる認証メール
	if request.TriggerSource == "CustomMessage_SignUp" || request.TriggerSource == "CustomMessage_ResendCode" {
		subscribeNews := false
		if sendSubscribeNews, ok := request.Request.ClientMetadata["subscribeNews"]; ok {
			if sendSubscribeNews == "1" {
				subscribeNews = true
			}
		}

		authenticationTokensCreator := domain.AuthenticationTokensCreator{
			CognitoSub:    request.UserName,
			SubscribeNews: subscribeNews,
			Time:          time.Now(),
		}

		scenario := application.CustomMessageScenario{
			Templates:                     templates,
			AuthenticationTokenRepository: authenticationTokenRepository,
			AuthenticationTokensCreator:   authenticationTokensCreator,
		}

		p := application.SignUpMessageBuildParams{Code: request.Request.CodeParameter, SubscribeNews: subscribeNews}
		body, err := scenario.BuildSignupMessage(p)
		if err != nil {
			// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
			log.Fatalln(err)
		}

		signupMessageResponse := events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body,
			EmailSubject: "サインアップ メールアドレスの確認をお願いします。",
		}

		request.Response = signupMessageResponse
	}

	if request.TriggerSource == "CustomMessage_ForgotPassword" {
		authenticationTokensCreator := domain.AuthenticationTokensCreator{
			CognitoSub: request.UserName,
			Time:       time.Now(),
		}

		scenario := application.CustomMessageScenario{
			Templates:                     templates,
			AuthenticationTokenRepository: authenticationTokenRepository,
			AuthenticationTokensCreator:   authenticationTokensCreator,
		}

		p := application.ForgotPasswordMessageBuildParams{Code: request.Request.CodeParameter}
		body, err := scenario.BuildForgotPasswordMessage(p)
		if err != nil {
			// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
			log.Fatalln(err)
		}

		forgotPasswordMessageResponse := events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body,
			EmailSubject: "パスワードをリセットします。",
		}

		request.Response = forgotPasswordMessageResponse
	}

	return request, nil
}

func main() {
	lambda.Start(handler)
}
