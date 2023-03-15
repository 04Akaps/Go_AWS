package main

import (
	"fmt"
	"log"

	"github.com/jjimgo/Go_AWS.git/goConfig"
	"github.com/jjimgo/Go_AWS.git/goaws"
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
	fmt.Println("Golang Aws Service Start")
}
