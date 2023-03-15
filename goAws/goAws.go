package goaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jjimgo/Go_AWS.git/errhandle"
)

func GetSession(regionPath string) *session.Session {
	// session은 재사용 가능하게 작성이 되어야 한다.

	// session, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(regionPath),
	// })

	// 이와 같이도 가능하지만 좀더 다양한 옵션 설정이 가능한 다음과 같은 방법을 추천

	awsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(regionPath),
		},
	})

	errhandle.ErrHandling(err)

	return awsSession
}
