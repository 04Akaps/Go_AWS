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
	// 이곳에 작성되는 추가적인 옵션이나 설명은 블로그에서 다룰 예정
	sendResult, err := goaws.sqsQueue.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: attributes,
		MessageBody:       aws.String(message),
		QueueUrl:          goaws.GetSQSQueuUrl(queueName),
	})
	errhandle.ErrHandling(err)

	fmt.Println(sendResult)
}

/*
메시지 수신 코드를 작성하기 전에 다음과 같은 사실을 알고 고려 해야 한다.

메시지를 수신하게 되면 어떤 처리를 해야 할지에 대해서 고민을 해야 한다.
SQS는 기본적으로 메시지를 수신한다고, 해당 메시지가 삭제 되는 것이 아니다.

그러기 떄문에 메시지의 목적에 따라서 해당 메시지를 삭제 할 것인가, 아니면 서비스에서도 수신이 되는지를 결정 해야 한다.
만약 한번 사용되고 처리 된다면, 그냥 삭제를 하면 된다.

문제는 또 발생한다.

마이크로서비스를 구축한다고 했을 떄 만약 A가 메시지를 수신했고, 삭제 처리를 시작했다고 하자.
하지만 이 떄 삭제가 이루어 지기 전에 B가 또 수신을 하게 된다면 프로세스가 꼬이게 된다.

이러한 시나리오를 피하고자 SQS는 가시성 타임아웃을 도입했다.
한 명의 소비자에 의해 수신된 메시지는 일정 시간동안 보이지 않게 되어서 다른 사용자는 해당 메시지를 볼 수 없게 한다.

또 주목할 점은 메시지를 항 상 두번 받지 않도록 보장을 할 수는 없다는 점이다.
왜냐하면 SQS큐들이 일반적을 다수의 서버에서 사용이 되기 떄문이다.

그러기 떄문에 오프라인이라 삭제 요청이 도달하지 못할 수도 있고, 메시지가 날아남을 가능성도 있다.

또다른 개념은 긴 폴링이나 대기시간이다.
분산되어 있기 떄문에 떄떄로 수신이 늦어 질 수 있다. 이러한 점을 고려한다면 들어오는 메시지에 대한 대기를 더 길게 가져가야 한다.
*/

// 이 함수는 특정 큐에 있는 메시지를 가져오는 함수가 된다.
// 테스트 용도로 작성이 되며, 나중에는 go루틴으로 돌릴 예정이다.
func (goaws *AwsSession) GetMessageFromSQS(queueName string) {
	// 이곳에 작성되는 추가적인 옵션이나 설명은 블로그에서 다룰 예정
	goaws.sqsQueue.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(*aws.String(sqs.MessageSystemAttributeNameSentTimestamp)),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            goaws.GetSQSQueuUrl(queueName),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(20),
	})
}
