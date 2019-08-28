package queuer

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	region          = "elasticmq"
	endpoint        = "http://localhost:9324"
	disableSSL      = true
	accessKeyID     = "x"
	secretAccessKey = "x"
	secretToken     = "x"
	queueNameTest   = "QueueURL"
)

// Queuer Holds the Queue Connection
type Queuer struct {
	svc *sqs.SQS
}

// New Initializes a Queuer
func New() *Queuer {
	return &Queuer{}
}

// Init the Queuer Structure to communicate with AWS SQS
func (q *Queuer) Init() error {
	//For simplicity, the Initialization will be done from hardcoded configuration
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
		DisableSSL:  aws.Bool(disableSSL),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, secretToken),
	}),
	)

	if sess == nil {
		return errors.New("Session not initialized properly")
	}

	// Create a SQS service client.
	q.svc = sqs.New(sess)

	if q.svc == nil {
		return errors.New("SQS not initialized properly")
	}

	return nil
}

// CreateQueue Creates a new Queue
func (q *Queuer) CreateQueue(name string) error {
	//For simplicity purposed, the Creation parameters will be done from hardcoded configuration
	_, err := q.svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(name),
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String("60"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	})
	return err
}

// GetQueueURL from QueueName
func (q *Queuer) GetQueueURL(queueName string) (string, error) {
	result, err := q.svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})

	if err != nil {
		return "", err
	}

	return *result.QueueUrl, nil
}

// WriteToQueue Writes a payload to a Queue
func (q *Queuer) WriteToQueue(queueURL string, body []byte) error {
	_, err := q.svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		// MessageAttributes: map[string]*sqs.MessageAttributeValue{
		// 	"Title": &sqs.MessageAttributeValue{
		// 		DataType:    aws.String("String"),
		// 		StringValue: aws.String("The Whistler"),
		// 	}
		// },
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(queueURL),
	})
	return err
}

// ReadFromQueue Reads a single value from queueUrl and deletes the message from Queue
func (q *Queuer) ReadFromQueue(queueURL string) (string, error) {
	result, err := q.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(20), // 20 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		return "", err
	}

	if len(result.Messages) == 0 {
		return "", fmt.Errorf("Received no messages from Queue %v", queueURL)
	}

	_, err = q.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: result.Messages[0].ReceiptHandle,
	})

	if err != nil {
		return "", err
	}

	//Since this is only receiving a single message from the queue
	return *result.Messages[0].Body, nil
}
