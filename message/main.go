package main

import (
	"bytes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/keitakn/go-cognito-lambda/infrastructure"
	"html/template"
	"log"
	"os"
)

var templates *template.Template

func init() {
	signupTemplatePath := "bin/signup-template.html"
	if infrastructure.IsTestRun() {
		currentDir, _ := os.Getwd()
		signupTemplatePath = currentDir + "/signup-template.html"
	}

	templates = template.Must(template.ParseFiles(signupTemplatePath))
}

type SignupMessage struct {
	ConfirmUrl string
}

func BuildSignupMessage(sm SignupMessage) (*bytes.Buffer, error) {
	var bodyBuffer bytes.Buffer

	err := templates.ExecuteTemplate(&bodyBuffer, "signup-template.html", sm)
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
		sm := SignupMessage{
			ConfirmUrl: "http://localhost:3900/cognito/signup/confirm?code=" + request.Request.CodeParameter + "&sub=" + request.UserName,
		}

		body, err := BuildSignupMessage(sm)
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
		forgotPasswordMessageResponse := events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage: "認証コードは {####} です。",
			EmailMessage: "次のリンクをクリックして、パスワードのリセットを完了させて下さい。 " +
				"http://localhost:3900/cognito/password/reset/confirm?code=" + request.Request.CodeParameter + "&sub=" + request.UserName,
			EmailSubject: "パスワードをリセットします。",
		}

		request.Response = forgotPasswordMessageResponse
	}

	return request, nil
}

func main() {
	lambda.Start(handler)
}
