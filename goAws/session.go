package goaws

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/jjimgo/Go_AWS/errhandle"
	"github.com/jjimgo/Go_AWS/goConfig"
)

type AwsSession struct {
	AwsSession *session.Session
	S3         *s3.S3
	S3Uploader *s3manager.Uploader
}

type CustomLogger struct {
	logFile *os.File
}

func (l *CustomLogger) Log(args ...interface{}) {
	log.SetOutput(l.logFile)
	log.Println(args...)
}

func GetSession(goConfig goConfig.GoConfig) *AwsSession {
	// session은 재사용 가능하게 작성이 되어야 한다.

	// session, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(regionPath),
	// })

	// 이와 같이도 가능하지만 좀더 다양한 옵션 설정이 가능한 다음과 같은 방법을 추천
	// 이렇게 .env보다는 os.GetEnv가 좀 더 개발하는데에 배포 서버를 관리하는데 편할 것라고 생각;;
	t := time.Now()
	startTime := t.Format("2006-01-02 15:04:05")
	logFile, err := os.Create("log/aws_" + startTime)
	errhandle.ErrHandling(err)

	logger := &CustomLogger{logFile: logFile}

	awsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(goConfig.AWS_REGION),
			Credentials: credentials.NewStaticCredentials(goConfig.IAM_ACCESS_KEY, goConfig.IAM_SECRET_KEY, ""),
			Logger:      logger,
			LogLevel:    aws.LogLevel(aws.LogDebugWithRequestErrors),
		},
	})
	errhandle.ErrHandling(err)

	svc := s3.New(awsSession)
	uploader := s3manager.NewUploader(awsSession)

	// result, err := uploader.Upload(uploadInput)

	return &AwsSession{
		AwsSession: awsSession,
		S3:         svc,
		S3Uploader: uploader,
	}
}
