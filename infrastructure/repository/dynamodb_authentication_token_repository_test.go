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

	if err := dynamodbHelper.CreateTestAuthenticationTokensTable(); err != nil {
		log.Fatal(err)
	}

	os.Exit(status)
}

func TestHandler(t *testing.T) {
	t.Run("Successful Create AuthenticationToken", func(t *testing.T) {
		repo := DynamodbAuthenticationTokenRepository{Dynamodb: db}

		token := "TestToken"

		putItem := domain.AuthenticationTokens{
			Token:          token,
			CognitoSub:     "0ef53af5-4eb9-4d2b-a939-8cb9d795512b",
			SubscribeNews:  true,
			ExpirationTime: 1922395084,
		}

		err := repo.Create(putItem)
		if err != nil {
			t.Fatal("Error failed to AuthenticationTokenRepository Create", err)
		}

		authenticationTokens, err := repo.FindByToken(token)
		if err != nil {
			t.Fatal("Error failed to AuthenticationTokenRepository FindByToken", err)
		}

		if reflect.DeepEqual(authenticationTokens, &putItem) == false {
			t.Error("\nActually: ", authenticationTokens, "\nExpected: ", putItem)
		}
	})
}
