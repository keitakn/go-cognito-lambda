package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ResponseBody struct {
	Message string `json:"message"`
}

// TODO 仮実装 後でちゃんとした実装にする
func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
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
