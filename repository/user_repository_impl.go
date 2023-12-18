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
	"log"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, db *helper.DynamoDB, user domain.User) domain.User {

	item, err := attributevalue.MarshalMap(user)
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
	return user
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, db *helper.DynamoDB, user domain.User) domain.User {
	userId, err := attributevalue.Marshal(user.Id)
	if err != nil {
		panic(err)
	}

	update := expression.Set(expression.Name("Name"), expression.Value(user.Name))
	update.Set(expression.Name("Email"), expression.Value(user.Email))

	if user.Password != "" {
		update.Set(expression.Name("Password"), expression.Value(user.Password))
	}

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
	} else {
		_, err := db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(db.TableName),
			Key:                       map[string]types.AttributeValue{"Id": userId},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			log.Printf("Couldn't update user. Here's why: %v\n", err)
		}
	}
	return user
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, db *helper.DynamoDB, user domain.User) {
	userId, err := attributevalue.Marshal(user.Id)
	if err != nil {
		panic(err)
	}
	_, err = db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(db.TableName),
		Key:       map[string]types.AttributeValue{"Id": userId},
	})
	if err != nil {
		panic(err)
	}
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, db *helper.DynamoDB, userId string) (domain.User, error) {
	user := domain.User{Id: userId}
	id, err := attributevalue.Marshal(user.Id)
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

	err = attributevalue.UnmarshalMap(response.Item, &user)
	if err != nil {
		panic(err)
	}
	return user, err
}
