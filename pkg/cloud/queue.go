package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Queue struct {
	client *sqs.Client
}

// NewQueue creates a new Queue object.
func NewQueue(cfg aws.Config) *Queue {
	return &Queue{
		client: sqs.NewFromConfig(cfg),
	}
}

// SendMessage sends a message to the specified SQS queue.
func (q *Queue) SendMessage(ctx context.Context, queueUrl string, messageBody string) (*string, error) {
	result, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueUrl),
		MessageBody: aws.String(messageBody),
	})
	if err != nil {
		return nil, err
	}
	return result.MessageId, nil
}

// ReceiveMessage receives messages from the specified SQS queue.
func (q *Queue) ReceiveMessage(ctx context.Context, queueUrl string, maxMessages int32) ([]types.Message, error) {
	result, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: maxMessages,
	})
	if err != nil {
		return nil, err
	}
	return result.Messages, nil
}

// DeleteMessage deletes a message from the specified SQS queue.
func (q *Queue) DeleteMessage(ctx context.Context, queueUrl string, receiptHandle string) error {
	_, err := q.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}
