package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"os"
)

type RequestBody struct {
	UserName    string `json:"userName"`
	NewPassword string `json:"newPassword"`
}

type ResponseBody struct {
	Message string `json:"message"`
}

var svc *cognitoidentityprovider.CognitoIdentityProvider

func init() {
	svc = cognitoidentityprovider.New(session.New(), &aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var reqBody RequestBody
	if err := json.Unmarshal([]byte(req.Body), &reqBody); err != nil {
		resBody := &ResponseBody{Message: "Bad Request"}
		resBodyJSON, _ := json.Marshal(resBody)

		res := events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:            string(resBodyJSON),
			IsBase64Encoded: false,
		}

		return res, err
	}

	param := &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(os.Getenv("TARGET_USER_POOL_ID")),
		Username:   aws.String(reqBody.UserName),
		Password:   aws.String(reqBody.NewPassword),
		Permanent:  aws.Bool(true),
	}

	if _, err := svc.AdminSetUserPassword(param); err != nil {
		resBody := &ResponseBody{Message: "failed to password update."}
		resBodyJSON, _ := json.Marshal(resBody)

		res := events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:            string(resBodyJSON),
			IsBase64Encoded: false,
		}

		return res, err
	}

	resBody := &ResponseBody{Message: "API Gateway v2 PATCH /users/passwords"}
	resBodyJSON, _ := json.Marshal(resBody)

	res := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(resBodyJSON),
		IsBase64Encoded: false,
	}

	return res, nil
}

func main() {
	lambda.Start(Handler)
}
