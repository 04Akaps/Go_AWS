package main

import (
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

func main() {
	// file, _ := os.Open("test.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }

	// awsSession.PutFileToS3("go-aws-test-bucket", "test.json", file)
	awsSession.GetFileFromS3("go-aws-test-bucket", "test.json")
	// awsSession.GetS3BucketList()
	// awsSession.GetAllObjectFromS3("go-aws-test-bucket")
}
