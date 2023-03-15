package goaws

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jjimgo/Go_AWS/errhandle"
)

func (goaws *AwsSession) getAclFromS3(bucket, fileKey string) {
	aclInput := &s3.GetObjectAclInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}

	result, err := goaws.S3.GetObjectAcl(aclInput)
	errhandle.ErrHandling(err)

	grants := result.Grants

	for _, grant := range grants {
		fmt.Println(*grant.Permission, *grant.Grantee.Type)
	}
}

func (goaws *AwsSession) PutFileToS3(fileName, fileKey string, file *os.File) {
	// 특정 객체를 업로드 하는 함수
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(fileName),
		Key:    aws.String(fileKey),
		Body:   file,
	}

	_, err := goaws.S3.PutObject(uploadInput)
	errhandle.ErrHandling(err)
}

type jsonTest struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func (goaws *AwsSession) GetFileFromS3(bucket, filekey string) []byte {
	// 특정 객체를 가져오는 함수
	goaws.getAclFromS3(bucket, filekey)
	downloadInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filekey),
	}

	res, err := goaws.S3.GetObject(downloadInput)
	errhandle.ErrHandling(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	errhandle.ErrHandling(err)

	return body
}

func (goaws *AwsSession) GetS3BucketList() {
	// 전체 bucket List를 가져오는 코드

	result, err := goaws.S3.ListBuckets(nil)
	errhandle.ErrHandling(err)

	for _, b := range result.Buckets {
		log.Printf("Bucket : %s \n", aws.StringValue(b.Name))
	}
}

func (goaws *AwsSession) GetAllObjectFromS3(bucket string) {
	// 버킷에 있는 모든 객체들을 가져오는 함수
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	err := goaws.S3.ListObjectsV2Pages(input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range page.Contents {
				fmt.Println("Object Key:", *obj.Key)
			}
			return !lastPage
		})

	errhandle.ErrHandling(err)
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
