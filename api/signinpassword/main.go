package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lestrrat-go/jwx/jwk"

	"github.com/keitakn/go-cognito-lambda/domain"
	"github.com/keitakn/go-cognito-lambda/infrastructure/repository"

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
	Password         string `json:"password"`
}

type IdTokenDetail struct {
	Jwt     string                `json:"jwt"`
	Payload domain.CognitoIdToken `json:"payload"`
}

type AccessTokenDetail struct {
	Jwt     string                    `json:"jwt"`
	Payload domain.CognitoAccessToken `json:"payload"`
}

type ResponseCreatedBody struct {
	IdToken      IdTokenDetail     `json:"idToken"`
	AccessToken  AccessTokenDetail `json:"accessToken"`
	RefreshToken string            `json:"refreshToken"`
}

type ResponseErrorBody struct {
	Message string `json:"message"`
}

var svc *cognitoidentityprovider.CognitoIdentityProvider
var iss string
var jwkSet jwk.Set

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

	iss = fmt.Sprintf(
		"https://cognito-idp.%v.amazonaws.com/%v",
		os.Getenv("REGION"),
		os.Getenv("TARGET_USER_POOL_ID"),
	)

	cognitoJwkRepository := &repository.HttpCognitoJwkRepository{
		Iss: iss,
	}

	res, err := cognitoJwkRepository.Fetch()
	if err != nil {
		// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
		log.Fatalln(err)
	}

	jwkSet = res
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

// TODO リファクタリングするまでの間、一時的に関数が大きい状態を許容する
//nolint:funlen
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

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(reqBody.Email),
			"PASSWORD": aws.String(reqBody.Password),
		},
		ClientId: aws.String(reqBody.UserPoolClientId),
	}

	resp, err := svc.InitiateAuth(input)
	if err != nil {
		// TODO 本来はライブラリのエラーメッセージをそのまま返してはいけない、適切なエラーメッセージに変換して返す事を推奨
		errorMessage := err.Error()

		resBody := &ResponseErrorBody{Message: errorMessage}
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(infrastructure.BadRequest, resBodyJson)

		return res, nil
	}

	cognitoJwtToken := &domain.CognitoJwtToken{
		Iss:    iss,
		Aud:    reqBody.UserPoolClientId,
		JwkSet: jwkSet,
	}

	idTokenPayload, err := cognitoJwtToken.ParseAndValidateIdToken(*resp.AuthenticationResult.IdToken)
	if err != nil {
		errorMessage := err.Error()

		resBody := &ResponseErrorBody{Message: errorMessage}
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(infrastructure.BadRequest, resBodyJson)

		return res, nil
	}

	idTokenDetail := &IdTokenDetail{
		Jwt:     *resp.AuthenticationResult.IdToken,
		Payload: *idTokenPayload,
	}

	accessTokenPayload, err := cognitoJwtToken.ParseAndValidateAccessToken(*resp.AuthenticationResult.AccessToken)
	if err != nil {
		errorMessage := err.Error()

		resBody := &ResponseErrorBody{Message: errorMessage}
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(infrastructure.BadRequest, resBodyJson)

		return res, nil
	}

	accessTokenDetail := &AccessTokenDetail{
		Jwt:     *resp.AuthenticationResult.AccessToken,
		Payload: *accessTokenPayload,
	}

	resBody := &ResponseCreatedBody{
		IdToken:      *idTokenDetail,
		AccessToken:  *accessTokenDetail,
		RefreshToken: *resp.AuthenticationResult.RefreshToken,
	}

	resBodyJson, _ := json.Marshal(resBody)

	res := createApiGatewayV2Response(infrastructure.Ok, resBodyJson)

	return res, nil
}

func main() {
	lambda.Start(Handler)
}
