package infrastructure

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

type DynamodbClientCreator struct{}

func (c *DynamodbClientCreator) Create() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}

func (c *DynamodbClientCreator) CreateTestClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess, &aws.Config{
		Endpoint:    aws.String(os.Getenv("DYNAMODB_TEST_ENDPOINT")),
		Region:      aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	})
}
