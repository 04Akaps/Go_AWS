package main

import (
	"log"
	"net/http"
	"os"
	"time"

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
	t := time.Now()
	startTime := t.Format("2006-01-02 15:04:05")
	logFile, err := os.Create("serverLog/server_" + startTime + ".log")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	errChannel := make(chan error)

	go func() {
		for {
			select {
			case err := <-errChannel:
				logger.Println(err)
			}
		}
	}()

	if err := http.ListenAndServe(":80", nil); err != nil {
		errChannel <- err
	}

	awsSession.SendMessageToSQS("golangEventQueue", "my Test Message")
}
