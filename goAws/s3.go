package goaws

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jjimgo/Go_AWS/errhandle"
)

func (goaws *AwsSession) PutFileToS3(fileName, fileKey string, file io.ReadSeeker) {
	file, err := os.Open("file.txt")
	errhandle.ErrHandling(err)
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(fileName),
		Key:    aws.String(fileKey),
		Body:   file,
	}

	_, err = goaws.S3.PutObject(uploadInput)
	errhandle.ErrHandling(err)
}

func (goaws *AwsSession) GetS3BucketList() {
	// 전체 bucket을 가져오는 코드

	result, err := goaws.S3.ListBuckets(nil)
	errhandle.ErrHandling(err)

	for _, b := range result.Buckets {
		log.Printf("Bucket : %s \n", aws.StringValue(b.Name))
	}
}

func (goaws *AwsSession) GetS3BucketPagination(bucketName string, maxKeys int64) {
	// pagination을 적용하겨 가져오는 코드

	inputparams := &s3.ListObjectsInput{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int64(maxKeys),
	}
	pageNum := 0
	goaws.S3.ListObjectsPages(inputparams, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		pageNum++

		for _, value := range page.Contents {
			fmt.Println(*value.Key)
		}

		return pageNum < 3
	})
}
