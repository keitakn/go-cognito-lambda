package test

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamodbHelper struct {
	Dynamodb *dynamodb.DynamoDB
}

func (h *DynamodbHelper) CreateTestAuthenticationTokensTable() error {
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
		TableName: aws.String(strings.Title(os.Getenv("DEPLOY_STAGE")) + "AuthenticationTokens"),
	}

	if _, err := h.Dynamodb.CreateTable(createTableInput); err != nil {
		return err
	}

	return nil
}

func (h *DynamodbHelper) DeleteTestAuthenticationTokensTable() error {
	deleteTableInput := &dynamodb.DeleteTableInput{
		TableName: aws.String(strings.Title(os.Getenv("DEPLOY_STAGE")) + "AuthenticationTokens"),
	}

	if _, err := h.Dynamodb.DeleteTable(deleteTableInput); err != nil {
		return err
	}

	return nil
}
