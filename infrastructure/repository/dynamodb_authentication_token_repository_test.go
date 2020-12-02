package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/keitakn/go-cognito-lambda/domain"
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

	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Token"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("CognitoSub"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Token"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("AuthenticationTokensGlobalIndexCognitoSub"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("CognitoSub"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String("AuthenticationTokens"),
	}

	if _, err := db.CreateTable(createTableInput); err != nil {
		log.Fatal(err)
	}

	status := m.Run()

	deleteTableInput := &dynamodb.DeleteTableInput{TableName: aws.String("AuthenticationTokens")}

	if _, err := db.DeleteTable(deleteTableInput); err != nil {
		log.Fatal(err)
	}

	os.Exit(status)
}

func TestHandler(t *testing.T) {
	t.Run("Successful Create AuthenticationToken", func(t *testing.T) {
		repo := DynamoDbAuthenticationTokenRepository{db: db}

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
