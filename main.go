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
	// 서버 로그 처리에 대해서 고민하다가 aws와 같은 방식으로 적용한 서버에 대한 로그 처리 채널
	// aws폴더에 있는 커스텀하게 작성한 interface를 사용하면 더 깔끔하게 작성 가능

	// t := time.Now()
	// startTime := t.Format("2006-01-02 15:04:05")
	// logFile, err := os.Create("serverLog/server_" + startTime + ".log")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer logFile.Close()

	// logger := log.New(logFile, "", log.LstdFlags)

	// errChannel := make(chan error)

	// go func() {
	// 	for {
	// 		select {
	// 		case err := <-errChannel:
	// 			logger.Println(err)
	// 		}
	// 	}
	// }()

	// if err := http.ListenAndServe(":80", nil); err != nil {
	// 	errChannel <- err
	// }

}
