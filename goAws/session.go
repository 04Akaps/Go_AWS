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
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/jjimgo/Go_AWS/errhandle"
	"github.com/jjimgo/Go_AWS/goConfig"
)

type AwsSession struct {
	AwsSession *session.Session
	S3         *s3.S3
	S3Uploader *s3manager.Uploader
	sqsQueue   *sqs.SQS
}

type CustomLogger struct {
	logFile *os.File
}

func (l *CustomLogger) Log(args ...interface{}) {
	log.SetOutput(l.logFile)
	log.Println(args...)
}

func GetSession(goConfig goConfig.GoConfig) *AwsSession {
	t := time.Now()
	startTime := t.Format("2006-01-02 15:04:05")
	logFile, err := os.Create("log/aws_" + startTime + ".log")
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
	sqsSvc := sqs.New(awsSession)

	return &AwsSession{
		AwsSession: awsSession,
		S3:         svc,
		S3Uploader: uploader,
		sqsQueue:   sqsSvc,
	}
}
