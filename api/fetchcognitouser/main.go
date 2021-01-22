package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/keitakn/go-cognito-lambda/infrastructure"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type ResponseOkBody struct {
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

func Handler(
	ctx context.Context, req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	cognitoSub := ""
	if _, ok := req.PathParameters["cognitoSub"]; ok {
		cognitoSub = req.PathParameters["cognitoSub"]
	}

	userPoolId := os.Getenv("TARGET_USER_POOL_ID")
	inputAdminGetUser := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &userPoolId,
		Username:   &cognitoSub,
	}

	user, err := svc.AdminGetUser(inputAdminGetUser)
	if err != nil {
		statusCode := infrastructure.InternalServerError
		errorMessage := "Internal Server Error"

		if err.Error() == "UserNotFoundException: User does not exist." {
			statusCode = infrastructure.NotFound
			errorMessage = "Not Found"
		}

		resBody := &ResponseErrorBody{Message: errorMessage}
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(statusCode, resBodyJson)

		return res, nil
	}

	if !*user.Enabled {
		statusCode := infrastructure.BadRequest
		errorMessage := "Cognito User Status Not Enabled"

		resBody := &ResponseErrorBody{Message: errorMessage}
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(statusCode, resBodyJson)

		return res, nil
	}

	resBody := &ResponseOkBody{CognitoSub: *user.Username}
	resBodyJson, _ := json.Marshal(resBody)

	res := createApiGatewayV2Response(infrastructure.Ok, resBodyJson)

	return res, nil
}

func main() {
	lambda.Start(Handler)
}
