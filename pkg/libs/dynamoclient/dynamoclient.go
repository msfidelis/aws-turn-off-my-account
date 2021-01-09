package dynamoclient

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type DynamoDbClient struct {
	tableName string
}

func New(tableName string) *DynamoDbClient {
	return &DynamoDbClient{
		tableName: tableName,
	}
}

func (d *DynamoDbClient) Save(item interface{}) *dynamodb.PutItemOutput {

	av, errMarsh := dynamodbattribute.MarshalMap(item)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	if errMarsh != nil {
		fmt.Println("Error to marshalling new item:")
		fmt.Println(errMarsh.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(d.tableName),
	}

	response, errPut := svc.PutItem(input)

	if errPut != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(errPut.Error())
		os.Exit(1)
	}

	return response
}

func (d DynamoDbClient) Scan(expr expression.Expression) *dynamodb.ScanOutput {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(d.tableName),
	}

	result, err := svc.Scan(params)

	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	return result
}

func (d DynamoDbClient) UpdateItem(keyMap map[string]*dynamodb.AttributeValue, expr expression.Expression) *dynamodb.UpdateItemOutput {

	fmt.Println(keyMap)
	fmt.Println(expr)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key:                       keyMap,
		ReturnValues:              aws.String("UPDATED_NEW"),
		TableName:                 aws.String(d.tableName),
		UpdateExpression:          expr.Update(),
	}

	result, err := svc.UpdateItem(input)

	if err != nil {
		fmt.Println("Updated API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	return result
}

func (d DynamoDbClient) RemoveItem(keyMap map[string]*dynamodb.AttributeValue) bool {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		Key:       keyMap,
		TableName: aws.String(d.tableName),
	}

	_, errDelete := svc.DeleteItem(input)

	if errDelete != nil {
		fmt.Println("Got error calling DeleteItem")
		fmt.Println((errDelete.Error()))
		return false
	} else {
		return true
	}

}
