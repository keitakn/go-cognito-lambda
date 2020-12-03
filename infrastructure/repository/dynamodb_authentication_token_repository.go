package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/keitakn/go-cognito-lambda/domain"
)

type DynamodbAuthenticationTokenRepository struct {
	db *dynamodb.DynamoDB
}

func (r *DynamodbAuthenticationTokenRepository) Create(item domain.AuthenticationTokens) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("AuthenticationTokens"),
		Item:      av,
	}

	_, err = r.db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (r *DynamodbAuthenticationTokenRepository) FindByToken(token string) (*domain.AuthenticationTokens, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("AuthenticationTokens"),
		Key: map[string]*dynamodb.AttributeValue{
			"Token": {
				S: aws.String(token),
			},
		},
	}

	result, err := r.db.GetItem(input)
	if err != nil {
		return nil, err
	}

	item := domain.AuthenticationTokens{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
