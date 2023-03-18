package goaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gopkg.in/mgo.v2/bson"
)

type DatabaseEvent struct {
	// DB에서 사용 될 데이터 구조
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name" dynamodbav:"EventName"` // 인덱스로 선언된 필드이기 떄문에
	Age       string        `bson:"age" json:"age"`
	Address   string        `bson:"address" json:"address"`
	DummyData []DummyEvent  `bson:"dummy_data" json:"dummy_data"`
}

type DummyEvent struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	DummyName string        `bson:"dummy_name" json:"dummy_name"`
	DummyAge  string        `bson:"dummy_age" json:"dummy_age"`
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

	if !event.ID.Valid() {
		event.ID = bson.NewObjectId()
	}

	for _, dummy := range event.DummyData {
		if !dummy.ID.Valid() {
			dummy.ID = bson.NewObjectId()
		}
	}

	newItem, err := dynamodbattribute.MarshalMap(event)

	if err != nil {
		return nil, err
	}

	_, err = dynamoDB.DynamoDBSession.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("events"),
		Item:      newItem,
	})

	if err != nil {
		return nil, err
	}

	return []byte(event.ID), nil
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
