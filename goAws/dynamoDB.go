package goaws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DatabaseEvent struct {
	// DB에서 사용 될 데이터 구조
	Name      string `json:"name"`
	Age       string `json:"age"`
	Address   string `json:"address"`
	DummyData []DummyEvent
}

type DummyEvent struct {
	DummyName string `json:"dummy_name"`
	DummyAge  string `json:"dummy_age"`
}

type DynamoDBLayout struct {
	DynamoDBSession *dynamodb.DynamoDB
}

type DatabaseHandler interface {
	// DB에서 사용할 Handler
	AddEvent(event DatabaseEvent) ([]byte, error)
	FindEvent([]byte) (DatabaseEvent, error)
	FindEventByName(string) (DatabaseEvent, error)
	FindAllEvents() ([]DatabaseEvent, error)
}

func NewDynamoDBClient(s *session.Session) (handler DatabaseHandler, err error) {
	if s == nil {
		s, err = session.NewSession()
		if err != nil {
			return nil, err
		}
	}

	dynamoDBClient := dynamodb.New(s)

	return &DynamoDBLayout{
		DynamoDBSession: dynamoDBClient,
	}, err
}

func (dynamoDB *DynamoDBLayout) AddEvent(event DatabaseEvent) ([]byte, error) {
	return nil, nil
}

func (dynamoDB *DynamoDBLayout) FindEvent(byte []byte) (DatabaseEvent, error) {
	return DatabaseEvent{}, nil
}

func (dynamoDB *DynamoDBLayout) FindEventByName(name string) (DatabaseEvent, error) {
	return DatabaseEvent{}, nil
}

func (dynamoDB *DynamoDBLayout) FindAllEvents() ([]DatabaseEvent, error) {
	return nil, nil
}
