package main

import (
	"flag"
	"github.com/Zhenhanyijiu/commongo/example/server"
	"github.com/Zhenhanyijiu/commongo/network"
	"github.com/google/uuid"
	"strconv"
	"sync"
)

type Param struct {
	uid           string
	guestNid      string
	hostNid       string
	networkType   string
	guestSendUrl  string
	hostSendUrl   string
	guestStartUrl string
	hostStartUrl  string
}

//./main -cmd=server -ip=127.0.0.2 -port=18002
//./main -cmd=sdk
func main() {
	cmds := flag.String("cmd", "server", "server cmd")
	ip := flag.String("ip", "127.0.0.1", "ip address")
	networkType := flag.String("net", "cache", "network type")
	port := flag.Int("port", 18001, "port")
	flag.Parse()
	uid := uuid.NewString()
	if *cmds == "sdk" {
		guestSendUrl := "http://127.0.0.1:18001/v1/send"
		hostSendUrl := "http://127.0.0.2:18002/v1/send"
		guestStartUrl := "http://127.0.0.1:18001/v1/algo/start"
		hosttStartUrl := "http://127.0.0.2:18002/v1/algo/start"
		param := Param{
			uid:           uid,
			guestNid:      "nid1",
			hostNid:       "nid2",
			networkType:   *networkType,
			guestSendUrl:  guestSendUrl,
			hostSendUrl:   hostSendUrl,
			guestStartUrl: guestStartUrl,
			hostStartUrl:  hosttStartUrl,
		}
		clientCmdSdk(param)
	}
	if *cmds == "server" {
		ser, err := server.NewServer(*ip+":"+strconv.Itoa(*port), *networkType)
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
		NetworkType: network.NetworkMap[param.networkType],
		Role:        server.GUEST,
		SendUrl:     param.hostSendUrl,
		Uid:         param.uid,
		LocalNid:    param.guestNid,
		RemoteNid:   param.hostNid,
	}
	go func() {
		defer wg.Done()
		server.StartTestSdk(&reqGuest, param.guestStartUrl)
	}()
	reqHost := server.Req{
		NetworkType: network.NetworkMap[param.networkType],
		Role:        server.HOST,
		SendUrl:     param.guestSendUrl,
		Uid:         param.uid,
		LocalNid:    param.hostNid,
		RemoteNid:   param.guestNid,
	}
	go func() {
		defer wg.Done()
		server.StartTestSdk(&reqHost, param.hostStartUrl)
	}()
	wg.Wait()
}
