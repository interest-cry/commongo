package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/Zhenhanyijiu/commongo/httpnetgo"
	"github.com/gin-gonic/gin"
)

var appwebnames = []string{"algo_0", "algo_1"}

func main_1() {
	httpnetgo.LocalDebug = true
	fmt.Printf(">>>>>>>localdebug:%+v\n", httpnetgo.LocalDebug)
	role := flag.String("r", "ser", "role,ser,cli")
	flag.Parse()
	cycN := 100000
	if *role == "ser" {
		run_server_1(cycN)
	} else {
		run_client_1(cycN)
	}
	time.Sleep(1 * time.Second)
}
func run_server_1(cycN int) {
	n := httpnetgo.NewNetServer("127.0.0.1:18888", "server001", []string{"algo_0", "algo_1"})
	//time.Sleep()
	con, err := n.NewServerConn(appwebnames[0], "1")
	fmt.Printf("con:%+v,err:%+v\n", con, err)
	if err != nil {
		return
	}
	defer func() {
		con.Close()
	}()
	for i := 0; i < cycN; i++ {
		ret, err := con.Recv()
		if err != nil {
			fmt.Printf("error 1:%+v\n", err)
			return
		}
		fmt.Printf("ret 1:%+v\n", string(ret))
		data := "server-data--" + strconv.Itoa(i)
		if err := con.Send([]byte(data)); err != nil {
			fmt.Printf("error 2:%+v\n", err)
			return
		}
	}
}
func run_client_1(cycN int) {
	n := httpnetgo.NewNetClient("127.0.0.2:28888")
	con, err := n.NewClientConn("127.0.0.1:18888", appwebnames[0])
	fmt.Printf("con:%+v,err:%+v\n", con, err)
	if err != nil {
		return
	}
	defer func() {
		con.Close()
	}()
	for i := 0; i < cycN; i++ {
		data := "client-data--" + strconv.Itoa(i)
		if err := con.Send([]byte(data)); err != nil {
			fmt.Printf("error 1:%+v\n", err)
			return
		}
		ret, err := con.Recv()
		if err != nil {
			fmt.Printf("error 2:%+v\n", err)
			return
		}
		fmt.Printf("ret:%+v\n", string(ret))
	}
}

func xx() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		_ = <-sigc
		fmt.Println("ctrl+c pressed")
		// client.logout()
	}()
}

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		_ = <-sigc
		fmt.Println(">>>>>>>>>> ctrl+c pressed")
	}()
	httpnetgo.LocalDebug = true
	fmt.Printf(">>>>>>>localdebug:%+v\n", httpnetgo.LocalDebug)
	role := flag.String("r", "ser", "role,ser,cli")
	flag.Parse()
	cycN := 10000
	if *role == "ser" {
		run_server_gin()
	} else {
		run_multi_thread(cycN)
	}
}

// type JobInfoType struct {
// 	Appname string
// 	Sessid  string
// }

func run_server_gin() {
	n := httpnetgo.NewNetServer("127.0.0.1:18888", "server001", []string{"algo_0", "algo_1"})
	eng := gin.New()
	eng.GET("/algo_0", func(c *gin.Context) {
		// var jobInfo JobInfoType
		// c.ShouldBindHeader(&jobInfo)
		appname := c.Request.Header.Get(httpnetgo.APPNAME)
		sessid := c.Request.Header.Get(httpnetgo.SESSID)
		// fmt.Printf("===== sessid:%+v,appname:%+v\n", sessid, appname)
		// fmt.Printf("client,FullPath:%v\n", c.FullPath())
		// fmt.Printf("client,URL:%v\n", c.Request.RequestURI)
		go func() {
			conn, err := n.NewServerConn(appname, sessid)
			if err != nil {
				panic("====== new server conn error")
			}
			defer func() { conn.Close() }()
			// conn := httpnet.NewConnect(jobInfo.AlgoId, jobInfo.TaskId,
			// 	jobInfo.Sid, n, jobInfo.ClientIp, jobInfo.ClientPort)
			for i := 0; i < 1; i++ {
				ret, err := conn.Recv()
				if err != nil {
					panic(err)
				}
				// fmt.Printf("sessid:%+v, i:%v ,ret:%+v\n", sessid, i, string(ret))
				fmt.Printf("sessid:%+v, i:%v ,ret:%+v\n", sessid, i, len(ret))
				// data := []byte(sessid + "#" + string(ret))
				data := ret
				if err := conn.Send(data); err != nil {
					panic(err)
				}
			}
		}()
	})
	fmt.Printf("start server process ...... 1\n")
	if err := eng.Run("127.0.0.1:9999"); err != nil {
		panic(err)
	}
}
func run_multi_thread(cycN int) {
	thnum := 10
	n := httpnetgo.NewNetClient("127.0.0.2:28888")
	wg := sync.WaitGroup{}
	wg.Add(thnum)
	for i := 0; i < thnum; i++ {
		go func(thid int) {
			defer func() { wg.Done() }()
			for ii := 0; ii < cycN; ii++ {
				run_client_algo_gin(n, thid)
			}
		}(i)
	}
	wg.Wait()
}
func run_client_algo_gin(n *httpnetgo.NetClient, thid int) {
	// n := httpnet.NewNet("127.0.0.2", 8888)
	// for i := 0; i < cyc; i++ {

	// taskId := "2000"
	// sid := "1000_" + strconv.Itoa(ind)
	// bjosn, _ := json.Marshal(&jobInfo)
	// go func() {
	// rsp, err := cli.Post("http://"+serverIp+":9999"+"/"+"algo_0", "application/json",
	// 	bytes.NewReader(bjosn))
	// serverIp := "127.0.0.1"
	// rq, err := http.NewRequest("GET", "http://"+serverIp+":9999"+"/"+"algo_0", nil)
	// if err != nil {
	// 	fmt.Printf("====== start req error:%v\n", err)
	// 	return
	// }
	// rq.Header.Set(httpnetgo.APPNAME, "algo_0")
	// rq.Header.Set(httpnetgo.SESSID)
	// rsp.Body.Close()
	// }()
	// n := httpnetgo.NewNetClient("127.0.0.2:28888")
	con, err := n.NewClientConn("127.0.0.1:18888", "algo_0")
	// fmt.Printf("con:%+v,err:%+v\n", con, err)
	if err != nil {
		return
	}
	defer func() {
		con.Close()
	}()
	con.SendStart("127.0.0.1:9999", "algo_0")
	// value := "value->" + sid
	maxLen := 1024 * 1024 * 20
	value := make([]byte, maxLen)
	// fmt.Printf("=========================1\n")
	// conn.Send(string(value))
	// fmt.Printf("=========================2\n")

	// value = "value----->" + sid
	for i := 0; i < 1; i++ {
		// data := fmt.Sprintf("[thid:%v]client-data--[sid:%v]", thid, con.GetSessid())
		// data := "client-data--" + con.GetSessid()
		data := value
		if err := con.Send(data); err != nil {
			fmt.Printf("error 1:%+v\n", err)
			return
		}
		ret, err := con.Recv()
		if err != nil {
			fmt.Printf("error 2:%+v\n", err)
			return
		}
		fmt.Printf("ret:%+v\n", len(ret))
	}
	// fmt.Printf("=========================ind:%v\n", ind)
}
