package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
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
		now := time.Now()

		expirationTime := now.Add(10 * time.Minute)

		repo := DynamodbAuthenticationTokenRepository{Dynamodb: db}

		token, err := uuid.NewRandom()
		if err != nil {
			t.Fatal("Error failed to Generate UUID", err)
		}

		putItem := domain.AuthenticationTokens{
			Token:          token.String(),
			CognitoSub:     "0ef53af5-4eb9-4d2b-a939-8cb9d795512b",
			SubscribeNews:  true,
			ExpirationTime: expirationTime.Unix(),
		}

		if err := repo.Create(putItem); err != nil {
			t.Fatal("Error failed to AuthenticationTokenRepository Create", err)
		}

		authenticationTokens, err := repo.FindByToken(token.String())
		if err != nil {
			t.Fatal("Error failed to AuthenticationTokenRepository FindByToken", err)
		}

		if reflect.DeepEqual(authenticationTokens, &putItem) == false {
			t.Error("\nActually: ", authenticationTokens, "\nExpected: ", putItem)
		}
	})
}
