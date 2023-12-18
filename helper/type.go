package helper

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

// DynamoDB represents a configuration for interacting with a DynamoDB database.
type DynamoDB struct {

	// Client is a pointer to a DynamoDB client that can be used
	// to perform database operation.
	// required
	Client *dynamodb.Client

	// TableName is the name of the DynamoDB table associated with
	// the instance.
	TableName string
}
