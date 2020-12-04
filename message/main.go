package main

import (
	"bytes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/keitakn/go-cognito-lambda/domain"
	"github.com/keitakn/go-cognito-lambda/infrastructure"
	"github.com/keitakn/go-cognito-lambda/infrastructure/repository"
	"html/template"
	"log"
	"os"
	"time"
)

var templates *template.Template
var db *dynamodb.DynamoDB
var authenticationTokenRepository domain.AuthenticationTokenRepository

func init() {
	signupTemplatePath := "bin/signup-template.html"
	forgotPasswordTemplatePath := "bin/forgot-password-template.html"

	dynamodbClientCreator := infrastructure.DynamodbClientCreator{}

	if infrastructure.IsTestRun() {
		currentDir, _ := os.Getwd()
		signupTemplatePath = currentDir + "/signup-template.html"
		forgotPasswordTemplatePath = currentDir + "/forgot-password-template.html"

		db = dynamodbClientCreator.CreateTestClient()
	}

	templates = template.Must(template.ParseFiles(signupTemplatePath, forgotPasswordTemplatePath))

	db = dynamodbClientCreator.Create()

	authenticationTokenRepository = &repository.DynamodbAuthenticationTokenRepository{Dynamodb: db}
}

type SignUpMessage struct {
	ConfirmUrl string
}

type ForgotPasswordMessage struct {
	ConfirmUrl string
}

func BuildSignupMessage(m SignUpMessage) (*bytes.Buffer, error) {
	var bodyBuffer bytes.Buffer

	err := templates.ExecuteTemplate(&bodyBuffer, "signup-template.html", m)
	if err != nil {
		return nil, err
	}

	return &bodyBuffer, nil
}

func BuildForgotPasswordMessage(m ForgotPasswordMessage) (*bytes.Buffer, error) {
	var bodyBuffer bytes.Buffer

	err := templates.ExecuteTemplate(&bodyBuffer, "forgot-password-template.html", m)
	if err != nil {
		return nil, err
	}

	return &bodyBuffer, nil
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

		authenticationTokens, err := authenticationTokensCreator.Create()
		if err != nil {
			// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
			log.Fatalln(err)
		}

		err = authenticationTokenRepository.Create(*authenticationTokens)
		if err != nil {
			// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
			log.Fatalln(err)
		}

		m := SignUpMessage{
			ConfirmUrl: "http://localhost:3900/cognito/signup/confirm?code=" + request.Request.CodeParameter + "&sub=" + request.UserName + "&authenticationToken=" + authenticationTokens.Token,
		}

		body, err := BuildSignupMessage(m)
		if err != nil {
			// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
			log.Fatalln(err)
		}

		signupMessageResponse := events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body.String(),
			EmailSubject: "サインアップ メールアドレスの確認をお願いします。",
		}

		request.Response = signupMessageResponse
	}

	if request.TriggerSource == "CustomMessage_ForgotPassword" {
		m := ForgotPasswordMessage{
			ConfirmUrl: "http://localhost:3900/cognito/password/reset/confirm?code=" + request.Request.CodeParameter + "&sub=" + request.UserName,
		}

		body, err := BuildForgotPasswordMessage(m)
		if err != nil {
			// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
			log.Fatalln(err)
		}

		forgotPasswordMessageResponse := events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body.String(),
			EmailSubject: "パスワードをリセットします。",
		}

		request.Response = forgotPasswordMessageResponse
	}

	return request, nil
}

func main() {
	lambda.Start(handler)
}
