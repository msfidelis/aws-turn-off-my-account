package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"serverless-architecture-boilerplate-go/pkg/libs/dynamoclient"
	"serverless-architecture-boilerplate-go/pkg/libs/sqsclient"
	"serverless-architecture-boilerplate-go/pkg/models/book"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {
		return errors.New("No SQS message passed to function")
	}

	dynamoTable := os.Getenv("DYNAMO_TABLE_BOOKS")
	dynamo := dynamoclient.New(dynamoTable)

	sqsQueue := os.Getenv("SQS_QUEUE_BOOKS")
	sqs := sqsclient.New(sqsQueue)

	for _, msg := range sqsEvent.Records {

		fmt.Printf("Got SQS message %q with body %q\n", msg.MessageId, msg.Body)
		book := &book.Book{}
		json.Unmarshal([]byte(msg.Body), book)
		book.Processed = true

		key := map[string]*dynamodb.AttributeValue{
			"hashkey": {
				S: aws.String(book.Hashkey),
			},
		}
		update := expression.Set(
			expression.Name("updated"),
			expression.Value(time.Now().String()),
		)
		update.Set(
			expression.Name("processed"),
			expression.Value(book.Processed),
		)
		expr, err := expression.NewBuilder().WithUpdate(update).Build()
		if err != nil {
			fmt.Printf("Caught Builder %v\n", err)
		}
		dynamo.UpdateItem(key, expr)
		sqs.DeleteMessage(msg.ReceiptHandle)
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
