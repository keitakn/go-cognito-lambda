package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/keitakn/go-cognito-lambda/infrastructure"
)

type RequestBody struct {
	UserPoolClientId string `json:"userPoolClientId"`
	ConfirmationCode string `json:"confirmationCode"`
	CognitoSub       string `json:"cognitoSub"`
	NewPassword      string `json:"newPassword"`
}

type ResponseErrorBody struct {
	Message string `json:"message"`
}

var svc *cognitoidentityprovider.CognitoIdentityProvider

//nolint:gochecknoinits
func init() {
	sess, err := session.NewSession()
	if err != nil {
		// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
		log.Fatalln(err)
	}

	svc = cognitoidentityprovider.New(sess, &aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
}

func createApiGatewayV2Response(statusCode int, resBodyJson []byte) events.APIGatewayV2HTTPResponse {
	res := events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(resBodyJson),
		IsBase64Encoded: false,
	}

	return res
}

func createApiGatewayV2NoContentResponse() events.APIGatewayV2HTTPResponse {
	res := events.APIGatewayV2HTTPResponse{
		StatusCode: infrastructure.NoContent,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		IsBase64Encoded: false,
	}

	return res
}

func Handler(
	ctx context.Context, req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	var reqBody RequestBody
	if err := json.Unmarshal([]byte(req.Body), &reqBody); err != nil {
		resBody := &ResponseErrorBody{Message: "Bad Request"}
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(infrastructure.BadRequest, resBodyJson)

		return res, err
	}

	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(reqBody.UserPoolClientId),
		Username:         aws.String(reqBody.CognitoSub),
		Password:         aws.String(reqBody.NewPassword),
		ConfirmationCode: aws.String(reqBody.ConfirmationCode),
	}

	_, err := svc.ConfirmForgotPassword(input)
	if err != nil {
		errorMessage := err.Error()

		switch errorMessage {
		case "CodeMismatchException: Invalid verification code provided, please try again.":
			errorMessage = "確認コードが一致しません、もう一度試して下さい。"
		case "ExpiredCodeException: Invalid code provided, please request a code again.":
			// 違うCognitoSubが指定された場合でもこのエラーメッセージが返ってくる
			errorMessage = "確認コードが無効、または有効期限切れです。"
		default:
			// TODO 本来はライブラリのエラーメッセージをそのまま返してはいけない、適切なエラーメッセージに変換して返す事を推奨
			errorMessage = err.Error()
		}

		resBody := &ResponseErrorBody{Message: errorMessage}
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(infrastructure.BadRequest, resBodyJson)

		return res, nil
	}

	res := createApiGatewayV2NoContentResponse()

	return res, nil
}

func main() {
	lambda.Start(Handler)
}
