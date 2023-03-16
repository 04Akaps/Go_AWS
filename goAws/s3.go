package goaws

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func makeS3Matadata(metadataKey, metadataValue string) map[string]*string {
	metadata := make(map[string]*string)

	metadata[metadataKey] = aws.String(metadataValue)
	return metadata
}

func (goaws *AwsSession) PutJsonFileToS3(fileName, fileKey, metadataKey, metadataValue string, byteData []byte) {
	err := ioutil.WriteFile("tempUploadFile", byteData, 0o644)
	errhandle.ErrHandling(err)
	defer os.Remove("tempUploadFile")

	file, err := os.Open("tempUploadFile")

	defer file.Close()

	goaws.PutFileToS3(fileName, fileKey, metadataKey, metadataValue, file)
}

func (goaws *AwsSession) PutFileToS3(fileName, fileKey, metadataKey, metadataValue string, file *os.File) {
	fileState, err := file.Stat()
	errhandle.ErrHandling(err)
	// 파일 사이즈에 따라서 다른 방식을 사용하는 것이 효율적이기 떄문에 분기를 쳐준다
	// Size의 반환값이 byte이기 떄문에 100MB는 다음과 같은 수치를 가진다.
	if fileState.Size() <= 100000000 {
		goaws.putFileToS3UsingPutObject(fileName, fileKey, metadataKey, metadataValue, file)
	} else {
		goaws.pubFileToS3UsingUploader(fileName, fileKey, metadataKey, metadataValue, file)
	}
}

func (goaws *AwsSession) putFileToS3UsingPutObject(fileName, fileKey, metadataKey, metadataValue string, file *os.File) {
	// 특정 객체를 업로드 하는 함수
	// 뭐든지 context를 사용하는 것이 좀 더 유연하고 좋기는 하다.
	// 내부적으로 효율적으로 돌아가게 구성이 되어 있으니
	uploadInput := &s3.PutObjectInput{
		Bucket:   aws.String(fileName),
		Key:      aws.String(fileKey),
		Body:     file,
		Metadata: makeS3Matadata(metadataKey, metadataValue),
	}

	_, err := goaws.S3.PutObject(uploadInput)
	errhandle.ErrHandling(err)
}

func (goaws *AwsSession) pubFileToS3UsingUploader(fileName, fileKey, metadataKey, metadataValue string, file *os.File) {
	// Uploader와 PutObject의 차이는 PutObject는 작은 용량의 파일을 업로드 할떄 유리
	// Uploader내부적으로 파일을 쪼개서 업로드 하기 떄문에 큰 용량의 파일을 업로드 할 떄 유리하다.
	// 대략적으로 100MB을 기준으로 사용 하면 된다.
	uploadInput := &s3manager.UploadInput{
		Bucket:   aws.String(fileName),
		Key:      aws.String(fileKey),
		Body:     file,
		Metadata: makeS3Matadata(metadataKey, metadataValue),
	}

	_, err := goaws.S3Uploader.Upload(uploadInput)
	errhandle.ErrHandling(err)
}

func (goaws *AwsSession) GetFileFromS3(bucket, filekey string) []byte {
	// 특정 객체를 가져오는 함수
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
	// pagination을 적용하여 bucket에 있는 객체들을 가져 오는 방법

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
