package cloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

var (
	// ErrEmptyQueueURL is returned when a required SQS queue URL is empty.
	ErrEmptyQueueURL = errors.New("empty queue URL")

	// ErrEmptyMessageBody is returned when the message body is missing in a send request.
	ErrEmptyMessageBody = errors.New("empty message body")

	// ErrEmptyReceiptHandle is returned when the receipt handle is not provided for delete or visibility operations.
	ErrEmptyReceiptHandle = errors.New("empty receipt handle")

	// ErrInvalidMaxMessages is returned when the number of requested messages is out of valid range (1 to 10).
	ErrInvalidMaxMessages = errors.New("max messages should be between 1 and 10")
)

const (
	defaultVisibilityTimeout = int32(30) // 30 seconds
	defaultWaitTimeSeconds   = int32(20) // 20 seconds for long polling
	maxMessagesLimit         = int32(10) // AWS SQS limit
)

// Queue represents an SQS client for queue operations.
type Queue struct {
	client *sqs.Client
}

// NewQueue creates a new instance of SQS client.
func NewQueue(cfg aws.Config) *Queue {
	return &Queue{
		client: sqs.NewFromConfig(cfg),
	}
}

// SendMessageInput represents input parameters for sending a message.
type SendMessageInput struct {
	QueueURL        string
	MessageBody     string
	DelaySeconds    *int32
	MessageGroupID  *string // For FIFO queues
	DeduplicationID *string // For FIFO queues
	Attributes      map[string]string
}

// ReceiveMessageInput represents input parameters for receiving messages.
type ReceiveMessageInput struct {
	QueueURL          string
	MaxMessages       int32
	VisibilityTimeout *int32
	WaitTimeSeconds   *int32
	AttributeNames    []types.QueueAttributeName
	MessageAttributes []string
}

// validate input parameters for SendMessage.
func (i *SendMessageInput) validate() error {
	if i.QueueURL == "" {
		return ErrEmptyQueueURL
	}
	if i.MessageBody == "" {
		return ErrEmptyMessageBody
	}
	return nil
}

// validate input parameters for ReceiveMessage.
func (i *ReceiveMessageInput) validate() error {
	if i.QueueURL == "" {
		return ErrEmptyQueueURL
	}
	if i.MaxMessages < 1 || i.MaxMessages > maxMessagesLimit {
		return ErrInvalidMaxMessages
	}
	return nil
}

// SendMessage sends a message to the specified SQS queue.
func (q *Queue) SendMessage(ctx context.Context, input SendMessageInput) (*string, error) {
	if err := input.validate(); err != nil {
		return nil, err
	}

	msgInput := &sqs.SendMessageInput{
		QueueUrl:    aws.String(input.QueueURL),
		MessageBody: aws.String(input.MessageBody),
	}

	if input.DelaySeconds != nil {
		msgInput.DelaySeconds = *input.DelaySeconds
	}
	if input.MessageGroupID != nil {
		msgInput.MessageGroupId = input.MessageGroupID
	}
	if input.DeduplicationID != nil {
		msgInput.MessageDeduplicationId = input.DeduplicationID
	}
	if len(input.Attributes) > 0 {
		msgInput.MessageAttributes = make(map[string]types.MessageAttributeValue)
		for k, v := range input.Attributes {
			msgInput.MessageAttributes[k] = types.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(v),
			}
		}
	}
	result, err := q.client.SendMessage(ctx, msgInput)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	return result.MessageId, nil
}

// ReceiveMessage receives messages from the specified SQS queue.
func (q *Queue) ReceiveMessage(ctx context.Context, input ReceiveMessageInput) ([]types.Message, error) {
	if err := input.validate(); err != nil {
		return nil, err
	}

	msgInput := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(input.QueueURL),
		MaxNumberOfMessages: input.MaxMessages,
		VisibilityTimeout:   defaultVisibilityTimeout,
		WaitTimeSeconds:     defaultWaitTimeSeconds,
	}

	if input.VisibilityTimeout != nil {
		msgInput.VisibilityTimeout = *input.VisibilityTimeout
	}

	if input.WaitTimeSeconds != nil {
		msgInput.WaitTimeSeconds = *input.WaitTimeSeconds
	}

	if len(input.AttributeNames) > 0 {
		msgInput.MessageSystemAttributeNames = make([]types.MessageSystemAttributeName, len(input.AttributeNames))
		for i, attr := range input.AttributeNames {
			msgInput.MessageSystemAttributeNames[i] = types.MessageSystemAttributeName(attr)
		}
	}

	if len(input.MessageAttributes) > 0 {
		msgInput.MessageAttributeNames = input.MessageAttributes
	}

	result, err := q.client.ReceiveMessage(ctx, msgInput)
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages: %w", err)
	}
	return result.Messages, nil
}

// DeleteMessage deletes a message from the specified SQS queue.
func (q *Queue) DeleteMessage(ctx context.Context, queueURL, receiptHandle string) error {
	if queueURL == "" {
		return ErrEmptyQueueURL
	}
	if receiptHandle == "" {
		return ErrEmptyReceiptHandle
	}

	_, err := q.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}

// ChangeMessageVisibility changes the visibility timeout of a message.
func (q *Queue) ChangeMessageVisibility(ctx context.Context, queueURL, receiptHandle string, timeoutSeconds int32) error {
	if queueURL == "" {
		return ErrEmptyQueueURL
	}
	if receiptHandle == "" {
		return ErrEmptyReceiptHandle
	}

	_, err := q.client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(queueURL),
		ReceiptHandle:     aws.String(receiptHandle),
		VisibilityTimeout: timeoutSeconds,
	})
	if err != nil {
		return fmt.Errorf("failed to change message visibility: %w", err)
	}
	return nil
}
