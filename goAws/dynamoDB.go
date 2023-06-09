package goaws

import (
	"errors"
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

func (dynamoLayout *DynamoDBLayout) AddEvent(event DatabaseEvent) ([]byte, error) {
	// event에 대한 항목은 서버측 controller에서 진행 할 부분임
	// 여기에서는 비지니스 로직만 담당

	// 아래 있는 두 코드도 controller에서 처리 해야 할 문제이기는 함
	//if !event.ID.Valid() {
	//	event.ID = bson.NewObjectId()
	//}
	//
	//for _, dummy := range event.DummyData {
	//	if !dummy.ID.Valid() {
	//		dummy.ID = bson.NewObjectId()
	//	}
	//}
	newItem, err := dynamodbattribute.MarshalMap(event)

	if err != nil {
		return nil, err
	}

	_, err = dynamoLayout.DynamoDBSession.PutItem(&dynamodb.PutItemInput{
		TableName:           aws.String("events"),
		Item:                newItem,
		ConditionExpression: aws.String("attribute_not_exists(id) OR id <>"),
		Expected: map[string]*dynamodb.ExpectedAttributeValue{
			"_id": {
				Value: &dynamodb.AttributeValue{N: aws.String(event.ID.Hex())},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return []byte(event.ID), nil
}

func (dynamoLayout *DynamoDBLayout) FindEvent(id []byte) (DatabaseEvent, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				B: id,
			},
		},
		TableName: aws.String("events"), // aws에서 테이블을 events로 선언
	}

	result, err := dynamoLayout.DynamoDBSession.GetItem(input)

	if err != nil {
		return DatabaseEvent{}, err
	}

	event := DatabaseEvent{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &event)

	return event, err
}

func (dynamoLayer *DynamoDBLayout) FindEventByName(name string) (DatabaseEvent, error) {
	input := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("EventName = :n"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {
				S: aws.String(name),
			},
		},
		IndexName: aws.String("EventName-index"),
		TableName: aws.String("events"),
	}

	result, err := dynamoLayer.DynamoDBSession.Query(input)

	if err != nil {
		return DatabaseEvent{}, nil
	}

	event := DatabaseEvent{}

	if len(result.Items) > 0 {
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &event)
	} else {
		err = errors.New("Not Found")
	}

	return event, err
}

func (dynamoLayout *DynamoDBLayout) FindAllEvents() ([]DatabaseEvent, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String("events"),
	}
	result, err := dynamoLayout.DynamoDBSession.Scan(input)

	if err != nil {
		return nil, err
	}

	events := []DatabaseEvent{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &events)

	return events, nil
}
