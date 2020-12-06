package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/keitakn/go-cognito-lambda/domain"
	"github.com/keitakn/go-cognito-lambda/test"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

var db *dynamodb.DynamoDB

func TestMain(m *testing.M) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	db = dynamodb.New(sess, &aws.Config{
		Endpoint:    aws.String("http://localhost:58000"),
		Region:      aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	})

	dynamodbHelper := test.DynamodbHelper{Dynamodb: db}

	if err := dynamodbHelper.CreateTestAuthenticationTokensTable(); err != nil {
		log.Fatal(err)
	}

	status := m.Run()

	if err := dynamodbHelper.DeleteTestAuthenticationTokensTable(); err != nil {
		log.Fatal(err)
	}

	os.Exit(status)
}

func TestHandler(t *testing.T) {
	t.Run("Successful Create AuthenticationToken", func(t *testing.T) {
		tokensCreator := domain.AuthenticationTokensCreator{
			CognitoSub:    "0ef53af5-4eb9-4d2b-a939-8cb9d795512b",
			SubscribeNews: true,
			Time:          time.Now(),
		}

		tokens, err := tokensCreator.Create()
		if err != nil {
			t.Fatal("Error failed to Generate AuthenticationTokens", err)
		}

		repo := DynamodbAuthenticationTokenRepository{Dynamodb: db}
		if err := repo.Create(*tokens); err != nil {
			t.Fatal("Error failed to AuthenticationTokenRepository Create", err)
		}

		authenticationTokens, err := repo.FindByToken(tokens.Token)
		if err != nil {
			t.Fatal("Error failed to AuthenticationTokenRepository FindByToken", err)
		}

		if reflect.DeepEqual(authenticationTokens, tokens) == false {
			t.Error("\nActually: ", authenticationTokens, "\nExpected: ", tokens)
		}
	})
}
