package sqsclient

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSClient struct {
	queueName string
}

func New(queueName string) *SQSClient {
	return &SQSClient{
		queueName: queueName,
	}
}

func (s *SQSClient) SendMessage(message interface{}) *sqs.SendMessageOutput {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sqs.New(sess)
	resultURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(s.queueName),
	})
	if err != nil {
		fmt.Println("Error to get QueueURL:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	m, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error to marshall struct")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	result, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(m)),
		QueueUrl:    resultURL.QueueUrl,
	})
	if err != nil {
		fmt.Println("Error to get send message to sqs:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return result
}

func (s *SQSClient) ReceiveMessage(MaxNumberOfMessages int64, VisibilityTimeout int64, WaitTimeSeconds int64) *sqs.ReceiveMessageOutput {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sqs.New(sess)
	resultURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(s.queueName),
	})
	if err != nil {
		fmt.Println("Error to get QueueURL:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	receive_params := &sqs.ReceiveMessageInput{
		QueueUrl:            resultURL.QueueUrl,
		MaxNumberOfMessages: aws.Int64(MaxNumberOfMessages),
		VisibilityTimeout:   aws.Int64(VisibilityTimeout),
		WaitTimeSeconds:     aws.Int64(WaitTimeSeconds),
	}
	receive_resp, err := svc.ReceiveMessage(receive_params)
	if err != nil {
		log.Println(err)
	}
	return receive_resp
}

func (s *SQSClient) DeleteMessage(receiptHandle string) bool {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sqs.New(sess)
	resultURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(s.queueName),
	})
	if err != nil {
		fmt.Println("Error to get QueueURL:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      resultURL.QueueUrl,
		ReceiptHandle: aws.String(receiptHandle),
	}
	_, errDel := svc.DeleteMessage(deleteParams)
	if errDel != nil {
		log.Println(err)
		os.Exit(1)
	}
	return true
}
