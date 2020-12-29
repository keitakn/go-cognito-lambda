package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/keitakn/go-cognito-lambda/domain"
	"github.com/keitakn/go-cognito-lambda/infrastructure"
	"github.com/keitakn/go-cognito-lambda/infrastructure/repository"
)

var db *dynamodb.DynamoDB
var authenticationTokenRepository domain.AuthenticationTokenRepository

func init() {
	dynamodbClientCreator := infrastructure.DynamodbClientCreator{}

	if infrastructure.IsTestRun() {
		db = dynamodbClientCreator.CreateTestClient()
	} else {
		db = dynamodbClientCreator.Create()
	}

	authenticationTokenRepository = &repository.DynamodbAuthenticationTokenRepository{Dynamodb: db}
}

func Handler(event events.CognitoEventUserPoolsVerifyAuthChallenge) (events.CognitoEventUserPoolsVerifyAuthChallenge, error) {
	targetUserPoolId := os.Getenv("TARGET_USER_POOL_ID")
	if targetUserPoolId != event.UserPoolID {
		return event, nil
	}

	event.Response.AnswerCorrect = false

	requestAuthenticationToken, ok := event.Request.ChallengeAnswer.(string)
	if ok == false {
		return event, nil
	}

	authenticationTokens, err := authenticationTokenRepository.FindByToken(requestAuthenticationToken)
	if err != nil {
		return event, err
	}

	userName := event.Request.PrivateChallengeParameters["answer"]
	if userName == authenticationTokens.CognitoSub {
		event.Response.AnswerCorrect = true
		return event, nil
	}

	return event, nil
}

func main() {
	lambda.Start(Handler)
}
