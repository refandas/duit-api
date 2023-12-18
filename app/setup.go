package app

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/refandas/duit-api/helper"
	"time"
)

// ListTables retrieves a list of table names from the specified DynamoDB instance.
func ListTables(ctx context.Context, db *helper.DynamoDB) []string {
	var tableNames []string
	tables, err := db.Client.ListTables(
		ctx,
		&dynamodb.ListTablesInput{},
	)
	if err != nil {
		panic(err)
	} else {
		tableNames = tables.TableNames
	}
	return tableNames
}

// TableExists checks whether a DynamoDB table exists in the specified
// DynamoDB instance.
func TableExists(ctx context.Context, db *helper.DynamoDB) bool {
	exists := false
	tableNames := ListTables(ctx, db)

	for _, table := range tableNames {
		if db.TableName == table {
			exists = true
		}
	}
	return exists
}

// SetupClient creates and returns a new DynamoDB client based on the
// default setup of AWS.
func SetupClient(ctx context.Context) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	return client
}

// CreateTableUser creates a new DynamoDB table named `Users` for storing user data
// using the specified DynamoDB instance.
func CreateTableUser(ctx context.Context, db *helper.DynamoDB) error {
	_, err := db.Client.CreateTable(
		ctx,
		&dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("Id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("Id"),
					KeyType:       types.KeyTypeHash,
				},
			},
			TableName: aws.String(db.TableName),
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		},
	)
	if err != nil {
		panic(err)
	}

	waiter := dynamodb.NewTableExistsWaiter(db.Client)
	err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(db.TableName),
	}, 5*time.Minute)

	if err != nil {
		panic(err)
	}

	return err
}

// CreateTableSpending creates a new DynamoDB table named `Spending` for
// storing user's spending data using the specified DynamoDB instance.
//
// The `Spending` table has a hash key of `Id` and a Global Secondary Index (GSI)
// named `UserIndex` with a hash key of `UserId` and sort key of `Date`.
func CreateTableSpending(ctx context.Context, db *helper.DynamoDB) error {
	_, err := db.Client.CreateTable(
		ctx,
		&dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("Id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("UserId"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("Date"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("Id"),
					KeyType:       types.KeyTypeHash,
				},
			},
			GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("UserIndex"),
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("UserId"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: aws.String("Date"),
							KeyType:       types.KeyTypeRange,
						},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(3),
						WriteCapacityUnits: aws.Int64(1),
					},
				},
			},
			TableName: aws.String(db.TableName),
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(3),
				WriteCapacityUnits: aws.Int64(1),
			},
		},
	)
	if err != nil {
		panic(err)
	}

	waiter := dynamodb.NewTableExistsWaiter(db.Client)
	err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(db.TableName),
	}, 5*time.Minute)

	return err
}

// CreateTable creates new DynamoDB table using the specified creation  function
// and the provided DynamoDB instance.
func CreateTable(ctx context.Context, db *helper.DynamoDB, createTableFunc func(ctx2 context.Context, dynamoDB *helper.DynamoDB) error) {
	tableExists := TableExists(ctx, db)
	if !tableExists {
		err := createTableFunc(ctx, db)
		if err != nil {
			panic(err)
		}
	}
}

// DeleteTable deletes a DynamoDB table using the specified DynamoDB instance
func DeleteTable(ctx context.Context, db *helper.DynamoDB) error {
	if TableExists(ctx, db) {
		_, err := db.Client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String(db.TableName),
		})
		return err
	}
	return nil
}

// SetupDatabase sets up and returns a helper.DynamoDB instance with configured client
// and created tables for user data and spending data.
func SetupDatabase(ctx context.Context) helper.DynamoDB {
	client := SetupClient(ctx)
	db := helper.DynamoDB{Client: client}

	// Create the table "Users" for user data.
	db.TableName = "Users"
	CreateTable(ctx, &db, CreateTableUser)

	// Create the table "Spending" for user  spending data
	db.TableName = "Spending"
	CreateTable(ctx, &db, CreateTableSpending)

	fmt.Println("--- Setup Database Done")
	return db
}
