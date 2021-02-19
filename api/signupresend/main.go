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
	Email            string `json:"email"`
	Name             string `json:"name"`
}

type ResponseCreatedBody struct {
	CognitoSub string `json:"cognitoSub"`
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

func fetchCognitoUserByEmail(email string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	userPoolId := os.Getenv("TARGET_USER_POOL_ID")
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &userPoolId,
		Username:   &email,
	}

	user, err := svc.AdminGetUser(input)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func mightUpdateUserAttributes(reqBody RequestBody, user *cognitoidentityprovider.AdminGetUserOutput) error {
	if !*user.Enabled {
		return nil
	}

	// 認証がまだの場合だけユーザー情報をアップデートする
	if *user.UserStatus == "UNCONFIRMED" {
		updateUserAttributesInput := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
			UserPoolId: aws.String(os.Getenv("TARGET_USER_POOL_ID")),
			Username:   aws.String(reqBody.Email),
			UserAttributes: []*cognitoidentityprovider.AttributeType{
				{
					Name:  aws.String("name"),
					Value: aws.String(reqBody.Name),
				},
			},
		}

		if _, err := svc.AdminUpdateUserAttributes(updateUserAttributesInput); err != nil {
			return err
		}
	}

	return nil
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

	user, err := fetchCognitoUserByEmail(reqBody.Email)
	if err != nil {
		resBody := &ResponseErrorBody{Message: "User Not Found"}
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(infrastructure.BadRequest, resBodyJson)

		return res, nil
	}

	if err := mightUpdateUserAttributes(reqBody, user); err != nil {
		// この条件分岐に入る事はシステム障害等以外はあり得ないので、500エラーを返す
		resBody := &ResponseErrorBody{Message: "予期せぬエラーが発生しました。時間が経ってから再度お試し下さい。"}
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(infrastructure.InternalServerError, resBodyJson)

		return res, nil
	}

	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(reqBody.UserPoolClientId),
		Username: aws.String(reqBody.Email),
	}

	_, err = svc.ResendConfirmationCode(input)
	if err != nil {
		// TODO 本来はライブラリのエラーメッセージをそのまま返してはいけない、適切なエラーメッセージに変換して返す事を推奨
		errorMessage := err.Error()

		if errorMessage == "UsernameExistsException: An account with the given email already exists." {
			errorMessage = "そのemailは既に利用されています。必要に応じて認証メールの再送APIを実行して下さい。"
		}

		resBody := &ResponseErrorBody{Message: errorMessage}
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(infrastructure.BadRequest, resBodyJson)

		return res, nil
	}

	resBody := &ResponseCreatedBody{CognitoSub: *user.Username}
	resBodyJson, _ := json.Marshal(resBody)

	res := createApiGatewayV2Response(infrastructure.Ok, resBodyJson)

	return res, nil
}

func main() {
	lambda.Start(Handler)
}
