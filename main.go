package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jjimgo/Go_AWS.git/goConfig"
	"github.com/jjimgo/Go_AWS.git/goaws"
)

var (
	serverConfig goConfig.GoConfig
	awsSession   *session.Session
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	serverConfig = goConfig.LoadGoConfig(".")

	awsSession = goaws.GetSession(serverConfig.AWS_REGION)
}

func main() {
	fmt.Println("Golang Aws Service Start")
}
