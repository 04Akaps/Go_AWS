package goaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/jjimgo/Go_AWS.git/errhandle"
	"github.com/jjimgo/Go_AWS.git/goConfig"
)

type AwsSession struct {
	AwsSession *session.Session
}

func GetSession(goConfig goConfig.GoConfig) *AwsSession {
	// session은 재사용 가능하게 작성이 되어야 한다.

	// session, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(regionPath),
	// })

	// 이와 같이도 가능하지만 좀더 다양한 옵션 설정이 가능한 다음과 같은 방법을 추천
	// 이렇게 .env보다는 os.GetEnv가 좀 더 개발하는데에 배포 서버를 관리하는데 편할 것라고 생각;;

	awsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(goConfig.AWS_REGION),
			Credentials: credentials.NewStaticCredentials(goConfig.IAM_ACCESS_KEY, goConfig.IAM_SECRET_KEY, ""),
		},
	})

	errhandle.ErrHandling(err)

	return &AwsSession{
		AwsSession: awsSession,
	}
}
