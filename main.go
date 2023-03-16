package main

import (
	"encoding/json"
	"log"

	"github.com/jjimgo/Go_AWS/goConfig"
	"github.com/jjimgo/Go_AWS/goaws"
)

var (
	serverConfig goConfig.GoConfig
	awsSession   *goaws.AwsSession
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	serverConfig = goConfig.LoadGoConfig(".")

	awsSession = goaws.GetSession(serverConfig)
}

type Person struct {
	Name    string
	Age     int
	Address string
}

func main() {
	// file, _ := os.Open("file/test.json")
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// // scanner := bufio.NewScanner(file)
	// // for scanner.Scan() {
	// // 	fmt.Println(scanner.Text())
	// // }
	p := Person{"John", 30, "123 Main St"}
	jsonBytes, _ := json.Marshal(p)
	awsSession.PutJsonFileToS3("go-aws-test-bucket", "sdlksdlksd.json", "testMeatData", "abc", jsonBytes)
	// awsSession.PutFileToS3("go-aws-test-bucket", "test.json", "testMeatData", "abc", body)
	// awsSession.GetFileFromS3("go-aws-test-bucket", "제발 되라.json")
	// // awsSession.GetS3BucketList()
	// // awsSession.GetAllObjectFromS3("go-aws-test-bucket")
}
