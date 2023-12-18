package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/refandas/duit-api/exception"
	"github.com/refandas/duit-api/helper"
	"github.com/refandas/duit-api/model/domain"
	"time"
)

type SpendingRepositoryImpl struct {
}

func NewSpendingRepository() SpendingRepository {
	return &SpendingRepositoryImpl{}
}

func (repository *SpendingRepositoryImpl) Save(ctx context.Context, db *helper.DynamoDB, spending domain.Spending) domain.Spending {
	item, err := attributevalue.MarshalMap(spending)
	if err != nil {
		panic(err)
	}
	_, err = db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(db.TableName),
		Item:      item,
	})
	if err != nil {
		panic(err)
	}
	return spending
}

func (repository *SpendingRepositoryImpl) Update(ctx context.Context, db *helper.DynamoDB, spending domain.Spending) domain.Spending {
	spendingId, err := attributevalue.Marshal(spending.Id)
	if err != nil {
		panic(err)
	}

	update := expression.Set(expression.Name("Amount"), expression.Value(spending.Amount))
	update.Set(expression.Name("Date"), expression.Value(spending.Date))
	update.Set(expression.Name("Category"), expression.Value(spending.Category))
	update.Set(expression.Name("Title"), expression.Value(spending.Title))
	update.Set(expression.Name("Description"), expression.Value(spending.Description))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		panic(err)
	} else {
		_, err := db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(db.TableName),
			Key:                       map[string]types.AttributeValue{"Id": spendingId},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			panic(err)
		}
	}
	return spending
}

func (repository *SpendingRepositoryImpl) Delete(ctx context.Context, db *helper.DynamoDB, spending domain.Spending) {
	spendingId, err := attributevalue.Marshal(spending.Id)
	if err != nil {
		panic(err)
	}
	_, err = db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(db.TableName),
		Key:       map[string]types.AttributeValue{"Id": spendingId},
	})
	if err != nil {
		panic(err)
	}
}

func (repository *SpendingRepositoryImpl) FindById(ctx context.Context, db *helper.DynamoDB, spendingId string) (domain.Spending, error) {
	spending := domain.Spending{Id: spendingId}
	id, err := attributevalue.Marshal(spending.Id)
	if err != nil {
		panic(err)
	}

	response, err := db.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key:       map[string]types.AttributeValue{"Id": id},
	})

	if err != nil {
		panic(err)
	}
	if response.Item == nil {
		panic(exception.NewNotFoundError("item not found"))
	}

	err = attributevalue.UnmarshalMap(response.Item, &spending)
	if err != nil {
		panic(err)
	}
	return spending, err
}

func (repository *SpendingRepositoryImpl) FindByUserId(ctx context.Context, db *helper.DynamoDB, userId string) []domain.Spending {
	var spendings []domain.Spending

	keyExpression := expression.Key("UserId").Equal(expression.Value(userId)).
		And(expression.Key("Date").LessThanEqual(expression.Value(time.Now().UnixMilli())))

	expr, err := expression.NewBuilder().WithKeyCondition(keyExpression).Build()
	if err != nil {
		panic(err)
	}

	response, err := db.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(db.TableName),
		IndexName:                 aws.String("UserIndex"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ScanIndexForward:          aws.Bool(true),
	})

	if err != nil {
		panic(err)
	}
	if len(response.Items) == 0 {
		panic(exception.NewNotFoundError("user not found"))
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &spendings)
	if err != nil {
		panic(err)
	}

	return spendings
}
