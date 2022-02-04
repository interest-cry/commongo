package main

import (
	"flag"
	"github.com/Zhenhanyijiu/commongo/example/server"
	"strconv"
	"sync"
)

type Param struct {
	guestSendUrl  string
	hostSendUrl   string
	guestStartUrl string
	hosttStartUrl string
}

//./main -cmd=server -ip=127.0.0.2 -port=18002
//./main -cmd=sdk
func main() {
	cmds := flag.String("cmd", "server", "server cmd")
	ip := flag.String("ip", "127.0.0.1", "ip address")
	port := flag.Int("port", 18001, "port")
	flag.Parse()
	if *cmds == "sdk" {
		guestSendUrl := "http://127.0.0.1:18001/v1/send"
		hostSendUrl := "http://127.0.0.2:18002/v1/send"
		guestStartUrl := "http://127.0.0.1:18001/v1/algo/start"
		hosttStartUrl := "http://127.0.0.2:18002/v1/algo/start"
		param := Param{
			guestSendUrl:  guestSendUrl,
			hostSendUrl:   hostSendUrl,
			guestStartUrl: guestStartUrl,
			hosttStartUrl: hosttStartUrl,
		}
		clientCmdSdk(param)
	}
	if *cmds == "server" {
		ser, err := server.NewServer(*ip + ":" + strconv.Itoa(*port))
		if err != nil {
			panic(err)
		}
		ser.Route()
	}
}
func clientCmdSdk(param Param) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	reqGuest := server.Req{
		Role:    server.GUEST,
		SendUrl: param.hostSendUrl,
	}
	go func() {
		defer wg.Done()
		server.StartTestSdk(&reqGuest, param.guestStartUrl)
	}()
	reqHost := server.Req{
		Role:    server.HOST,
		SendUrl: param.guestSendUrl,
	}
	go func() {
		defer wg.Done()
		server.StartTestSdk(&reqHost, param.hosttStartUrl)
	}()
	wg.Wait()
}
