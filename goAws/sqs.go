package goaws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jjimgo/Go_AWS/errhandle"
)

// Blog : https://medium.com/@sdl182975/golang%EC%97%90%EC%84%9C%EC%9D%98-aws-3-48dbe9f0c290

type EventEmitter interface {
	Emit(event Event) error
}

type EventListener interface {
	Listen(events ...string) (<-chan Event, <-chan error, error)
	ReceiveMessage(eventCh chan Event, errorCh chan error, events ...string)
}

type SqsEmitter struct {
	SqsSvc   *sqs.SQS
	QueueURL *string
}

type SqsListener struct {
	SqsSvc              *sqs.SQS
	QueueURL            *string
	maxNumberOfMessages int64
	waitTime            int64
	visibilityTimeOut   int64
}

type Event struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func NewSqsEventEmitter(s *session.Session, queueName string) (EventEmitter, error) {
	sqsSvc := sqs.New(s)
	outInfo, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return nil, err
	}

	return &SqsEmitter{
		SqsSvc:   sqsSvc,
		QueueURL: outInfo.QueueUrl,
	}, nil
}

func NewSqsEventListener(s *session.Session, queueName string, maxMsg, waitTime, visibiliryTimeOut int64) (listener EventListener, err error) {
	if s == nil {
		s, err = session.NewSession()
		if err != nil {
			return
		}
	}

	svc := sqs.New(s)

	outInfo, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return
	}

	listener = &SqsListener{
		SqsSvc:              svc,
		QueueURL:            outInfo.QueueUrl,
		maxNumberOfMessages: maxMsg,
		waitTime:            waitTime,
		visibilityTimeOut:   visibiliryTimeOut,
	}

	return
}

func (sqsListener *SqsListener) Listen(events ...string) (<-chan Event, <-chan error, error) {
	if sqsListener == nil {
		return nil, nil, errors.New("SQSListener is nil : Listen()")
	}

	eventCh := make(chan Event)
	errorCh := make(chan error)

	go func() {
		for {
			sqsListener.ReceiveMessage(eventCh, errorCh)
		}
	}()

	return eventCh, errorCh, nil
}

func (sqsListener *SqsListener) ReceiveMessage(eventCh chan Event, errorCh chan error, events ...string) {
	recvMsgResult, err := sqsListener.SqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            sqsListener.QueueURL,
		MaxNumberOfMessages: aws.Int64(sqsListener.maxNumberOfMessages),
		WaitTimeSeconds:     aws.Int64(sqsListener.waitTime),
		VisibilityTimeout:   aws.Int64(sqsListener.visibilityTimeOut),
	})

	if err != nil {
		errorCh <- err
	}

	bContinue := false

	for _, msg := range recvMsgResult.Messages {
		value, ok := msg.MessageAttributes["event_name"]
		if !ok {
			continue
		}

		eventName := aws.StringValue(value.StringValue)

		for _, event := range events {
			if strings.EqualFold(eventName, event) {
				bContinue = true
				break
			}
		}

		if !bContinue {
			continue
		}

		message := aws.StringValue(msg.Body)
		var event Event

		err = json.Unmarshal([]byte(message), &event)

		if err != nil {
			errorCh <- err
		}

		eventCh <- event

		sqsListener.SqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      sqsListener.QueueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})

		if err != nil {
			errorCh <- err
		}
	}

}

func (SqsEmitter *SqsEmitter) Emit(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	attributes := map[string]*sqs.MessageAttributeValue{
		"event_name": {
			DataType:    aws.String("String"),
			StringValue: aws.String(event.Name),
		},
	}
	_, err = SqsEmitter.SqsSvc.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: attributes,
		MessageBody:       aws.String(string(data)),
		QueueUrl:          SqsEmitter.QueueURL,
	})
	return err
}

// -> 여기 아래에 있는 내용은 정말 Sample 함수들
// 인터페이스를 구축하고 사용하는 코드는 위에 있는 코드를 참고

func (sqsEmitter *SqsEmitter) SendMessageToSQS(queueName, message string) {
	attributes := map[string]*sqs.MessageAttributeValue{
		"message_type": {
			DataType:    aws.String("String"),
			StringValue: aws.String("RESERVATION"),
		},
		"Count": {
			DataType:    aws.String("Number"),
			StringValue: aws.String("2"),
		},
		"Binary_Type": {
			DataType:    aws.String("Binary"),
			BinaryValue: []byte{0x00, 0x01, 0x02}, // [0,1,2]
		},
	}
	sendResult, err := sqsEmitter.SqsSvc.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: attributes,
		MessageBody:       aws.String(message),
		QueueUrl:          sqsEmitter.QueueURL,
	})
	errhandle.ErrHandling(err)

	fmt.Println(sendResult)
}

func (sqsEmitter *SqsEmitter) GetMessageFromSQS(queueName string) {
	QueueUrl := sqsEmitter.QueueURL

	receivedMessage, err := sqsEmitter.SqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		// 가져올 정보를 입력하게 된다.
		AttributeNames: []*string{
			aws.String(*aws.String(sqs.MessageSystemAttributeNameSentTimestamp)),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            QueueUrl,
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(20),
	})
	errhandle.ErrHandling(err)
	fmt.Println(receivedMessage)
	// 아래에 있는 코드는 수신된 메시지를 처리하는 코드이다.

	for i, msg := range receivedMessage.Messages {
		log.Println("NewMessage", i, *msg.Body)

		for key, value := range msg.MessageAttributes {
			log.Println("Message Attributes", key, aws.StringValue(value.StringValue))
		}

		for key, value := range msg.Attributes {
			log.Println("Attrobite", key, *value)
		}

		log.Println("------ Delete Message -------")

		deletedInfo, err := sqsEmitter.SqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      QueueUrl,
			ReceiptHandle: msg.ReceiptHandle,
		})

		errhandle.ErrHandling(err)

		log.Println("Message Deleted...!!")
		fmt.Println(deletedInfo)
	}

	// 수신을 하게 되면 전송할 떄의 MessageAttributes정보, Body값들 모두 나오게 된다.
}
