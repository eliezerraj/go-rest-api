package dynamodb

import (
	"log"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/go-rest-api/internal/adapter/contract"
	"github.com/go-rest-api/internal/model"
	"github.com/go-rest-api/internal/error"
)

type BalanceRepositoryDynamoDBImpl struct {
	client dynamodbiface.DynamoDBAPI
	table_name  *string
}

func NewBalanceRepositoryDynamoDB(	table_name string,
									app model.ManagerInfo) (contract.BalanceRepositoryAdapterPort, error) {

	log.Print("NewBalanceRepositoryDynamoDB -----",table_name) 
	log.Print("NewBalanceRepositoryDynamoDB -----",app.AwsEnv.Aws_region) 
	log.Print("NewBalanceRepositoryDynamoDB -----",app.AwsEnv.Aws_access_id) 
	log.Print("NewBalanceRepositoryDynamoDB -----",app.AwsEnv.Aws_access_secret) 

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(app.AwsEnv.Aws_region),
		Credentials: credentials.NewStaticCredentials(	app.AwsEnv.Aws_access_id,
														app.AwsEnv.Aws_access_secret, 
														""),},
	)
	if err != nil {
		return BalanceRepositoryDynamoDBImpl{}, erro.ErrOpenDatabase
	}

	return BalanceRepositoryDynamoDBImpl{
		client: dynamodb.New(sess),
		table_name: aws.String(table_name),
	}, nil
}

func (b BalanceRepositoryDynamoDBImpl) AddBalance(ctx context.Context, balance model.Balance) (model.Balance, error) {
	log.Print("AddBalance-----") 
	
	item, err := dynamodbattribute.MarshalMap(balance)
	if err != nil {
		log.Print("erro :", err) 
		return model.Balance{}, erro.ErrSaveDatabase
	}

	transactItems := []*dynamodb.TransactWriteItem{}
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{Put: &dynamodb.Put{
		TableName: b.table_name,
		Item:      item,
	}})

	transaction := &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}
	if err := transaction.Validate(); err != nil {
		log.Print("erro :", err) 
		return model.Balance{}, erro.ErrSaveDatabase
	}

	_, err = b.client.TransactWriteItemsWithContext(ctx, transaction)
	if err != nil {
		log.Print("erro :", err) 
		return model.Balance{}, erro.ErrSaveDatabase
	}

	return balance , nil
}

func (b BalanceRepositoryDynamoDBImpl) ListBalance(ctx context.Context) ([]model.Balance, error) {
	log.Print("List") 
	return []model.Balance{}, erro.ErrListNotAllowed
}

func (b BalanceRepositoryDynamoDBImpl) ListBalanceById(ctx context.Context, pk string, sk string) ([]model.Balance, error) {
	log.Print("ListBalanceById") 

	balance_id := pk
	account := sk

	var keyCond expression.KeyConditionBuilder

	keyCond = expression.KeyAnd(
		expression.Key("balance_id").Equal(expression.Value(balance_id)),
		expression.Key("account").BeginsWith(account),
	)
	expr, err := expression.NewBuilder().
							WithKeyCondition(keyCond).
							Build()
	if err != nil {
		return []model.Balance{}, erro.ErrNotFound
	}

	key := &dynamodb.QueryInput{
		TableName:                 b.table_name,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	fmt.Println("key => ", key)

	result, err := b.client.QueryWithContext(ctx, key)
	if err != nil {
		log.Print("erro :", err) 
		return []model.Balance{}, erro.ErrNotFound
	}

	fmt.Println("result => ", result)

	balances := []model.Balance{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &balances)
	fmt.Println("--------------------------------")
	fmt.Println("balances => ", balances)
    if err != nil {
		log.Print("erro :", err) 
		return []model.Balance{}, erro.ErrUnmarshal
    }

	if len(balances) == 0 {
		return []model.Balance{}, erro.ErrNotFound
	} else {
		return balances, nil
	}
}

func (b BalanceRepositoryDynamoDBImpl) GetBalance(ctx context.Context, pk string) (model.Balance, error) {
	log.Print("GetBalance") 

	balance_id := pk

	var keyCond expression.KeyConditionBuilder
	keyCond = expression.Key("balance_id").Equal(expression.Value(balance_id))

	expr, err := expression.NewBuilder().
							WithKeyCondition(keyCond).
							Build()
	if err != nil {
		log.Print("erro :", err) 
		return model.Balance{}, erro.ErrNotFound
	}

	key := &dynamodb.QueryInput{
		TableName:                 b.table_name,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	fmt.Println("key => ", key)

	result, err := b.client.QueryWithContext(ctx, key)
	if err != nil {
		log.Print("erro :", err) 
		return model.Balance{}, erro.ErrNotFound
	}

	fmt.Println("result => ", result)

	balances := []model.Balance{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &balances)
    if err != nil {
		log.Print("erro :", err) 
		return model.Balance{}, erro.ErrUnmarshal
    }

	if len(balances) == 0 {
		return model.Balance{}, erro.ErrNotFound
	} else {
		return balances[0], nil
	}
}
