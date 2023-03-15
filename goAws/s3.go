package goaws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jjimgo/Go_AWS.git/errhandle"
)

func GetS3BucketList(session *AwsSession) {
	// 전체 bucket을 가져오는 코드
	s3Svc := s3.New(session.AwsSession)

	result, err := s3Svc.ListBuckets(nil)
	errhandle.ErrHandling(err)

	for _, b := range result.Buckets {
		log.Printf("Bucket : %s \n", aws.StringValue(b.Name))
	}
}

func GetS3BucketPagination(session *AwsSession, bucketName string, maxKeys int64) {
	// pagination을 적용하겨 가져오는 코드
	s3Svc := s3.New(session.AwsSession)

	inputparams := &s3.ListObjectsInput{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int64(maxKeys),
	}
	pageNum := 0
	s3Svc.ListObjectsPages(inputparams, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		pageNum++

		for _, value := range page.Contents {
			fmt.Println(*value.Key)
		}

		return pageNum < 3
	})
}
