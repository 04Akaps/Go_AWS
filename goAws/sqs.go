package goaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jjimgo/Go_AWS/errhandle"
)

/*
큐로 메시지를 보내는 데에는 다음과 같은 기본적인 개념을 숙지하는 것이 좋다.

1. Message body

단순히 보내고자 하는 메시지를 의미한다. 즉 만약 hello SQS 라는 Text를 보내고 싶다면, 단순히 해당 Text가 Message Body가 된다.

2. Message attributes

메시지 속성은 구조화된 메타데이터들의 모임이라고 할 수 있다.
단순히 키 - 값 쌍들의 리스트로 전송이 가능하다.
하지만 단순하게 설정을 하는 속성은 아니다.
해당 속성에는 내가 보내는 메시지가 어떤 메시지인지 설명을 하는 항목이기 떄문에 중요하다.
세가지 주요 타입을 제공한다.
- 문자열, 숫자, 바이너리
여기서 바이너리는 압축 파일이나 이미지 같은 바이너리를 의미
*/

func (goaws *AwsSession) GetSQSQueuUrl(queunName string) *string {
	QuInfo, err := goaws.sqsQueue.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queunName),
	})
	errhandle.ErrHandling(err)

	return QuInfo.QueueUrl
}

func (goaws *AwsSession) SendMessageToSQS(queueName, message string) {
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

	sendResult, err := goaws.sqsQueue.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: attributes,
		MessageBody:       aws.String(message),
		QueueUrl:          goaws.GetSQSQueuUrl(queueName),
	})
	errhandle.ErrHandling(err)

	fmt.Println(sendResult)
}
