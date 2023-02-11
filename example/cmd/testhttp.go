package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/Zhenhanyijiu/commongo/httpnet"
	"github.com/gin-gonic/gin"
)

// var url = "http://127.0.0.1:8888/send"
var cyc = 1
var serverIp = "127.0.0.1"
var serverPort = 8888

var clientIp = "127.0.0.2"
var clientPort = 18888
var MaxLen = 1024 * 1024 * 10

type JobInfo struct {
	AlgoId     string
	TaskId     string
	Sid        string
	ServerIp   string
	ServerPort int
	ClientIp   string
	ClientPort int
}

func main4() {
	cli := http.Client{Transport: http.DefaultTransport}
	req, err := http.NewRequest("POST", "http://127.0.0.1:18888/hi", bytes.NewReader([]byte{}))
	if err != nil {
		panic(err)
	}
	req.Header.Add("h1", "h1_vv")
	rsp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		rsp.Body.Close()
	}()
	ret, err := io.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ret:%v\n", string(ret))
}

func main() {
	role := flag.String("r", "ser", "role,ser,cli")
	flag.Parse()
	if *role == "ser" {
		run_server()
	} else {
		go_num := 1
		n := httpnet.NewNet(clientIp, clientPort)
		wg := sync.WaitGroup{}
		wg.Add(go_num)
		for iii := 0; iii < go_num; iii++ {
			go func(thrId int) {
				base := (thrId + 1) * 100000
				defer func() {
					wg.Done()
				}()
				for i := 0; i < cyc; i++ {
					run_client_algo(n, base+i)
				}
			}(iii)
		}
		wg.Wait()
	}

}

func run_server() {
	n := httpnet.NewNet(serverIp, serverPort)
	eng := gin.New()
	eng.POST("/algo_0", func(c *gin.Context) {
		var jobInfo JobInfo
		c.ShouldBindJSON(&jobInfo)
		fmt.Printf("===== jobInfo:%+v\n", jobInfo)
		fmt.Printf("client,FullPath:%v\n", c.FullPath())
		fmt.Printf("client,URL:%v\n", c.Request.RequestURI)
		go func() {
			conn := httpnet.NewConnect(jobInfo.AlgoId, jobInfo.TaskId,
				jobInfo.Sid, n, jobInfo.ClientIp, jobInfo.ClientPort)
			for i := 0; i < 1; i++ {
				ret := conn.Recv()
				fmt.Printf("====i:%v ,ret:%+v\n", i, len(ret.Value))
			}
			// ret = conn.Recv()
			// fmt.Printf("====2: ,ret:%+v\n", len(ret.Value))
			conn.Close()
		}()
	})
	fmt.Printf("start server process ...... 1\n")
	if err := eng.Run("127.0.0.1:9999"); err != nil {
		panic(err)
	}

}
func run_client_algo(n *httpnet.Net, ind int) {
	// n := httpnet.NewNet("127.0.0.2", 8888)
	// for i := 0; i < cyc; i++ {
	cli := http.Client{Transport: http.DefaultTransport}
	taskId := "2000"
	sid := "1000_" + strconv.Itoa(ind)
	jobInfo := JobInfo{AlgoId: "algo_0", TaskId: taskId, Sid: sid,
		ServerIp: serverIp, ServerPort: serverPort,
		ClientIp: clientIp, ClientPort: clientPort}
	bjosn, _ := json.Marshal(&jobInfo)
	// go func() {
	rsp, err := cli.Post("http://"+serverIp+":9999"+"/"+"algo_0", "application/json",
		bytes.NewReader(bjosn))
	if err != nil {
		fmt.Printf("====== start req error:%v\n", err)
		return
	}
	rsp.Body.Close()
	// }()

	////////////
	conn := httpnet.NewConnect(jobInfo.AlgoId, jobInfo.TaskId,
		jobInfo.Sid, n, jobInfo.ServerIp, jobInfo.ServerPort)
	// value := "value->" + sid
	value := make([]byte, MaxLen)
	// fmt.Printf("=========================1\n")
	// conn.Send(string(value))
	// fmt.Printf("=========================2\n")

	// value = "value----->" + sid
	for i := 0; i < 1; i++ {
		conn.Send(string(value[:1024*1024*5]))
	}
	fmt.Printf("=========================ind:%v\n", ind)
}
