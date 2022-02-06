package network

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

/***
go test -v -run TestIp
go test -v -test.run=TestIp
go test -v ./*.go
go test -v ./*.go -test.run=.
go test -v ./*.go -test.bench=.
go test -v options_test.go -test.run TestIp
***/
/**
会出现未定义的情况, 这是因为定义在其他文件里, 需要加上定义的文件.
go test -v options.go options_test.go
go test -v ./*.go -test.run=. 全测
go test -v ./*.go -test.run=TestTcpConn_NewTcpConn
**/
//go test -v hello.go hello_test.go
//go test -v  -test.run="TestA*";加前缀
/***
测试覆盖率
go test -v -cover
go test -v -coverprofile=a.out -test.run="TestA*" # 把测试结果保存在 a.out
go tool cover -html=./a.out  # 通过浏览器打开, 可以看到覆盖经过的函数
 ****/

/*
benchmark测试
命令：
go test -bench=. -benchmem
指定特定的一个benchmark:
go test -v -bench=BenchmarkIp -benchmem -run=^$
go test -v -bench=BenchmarkIp -benchmem -run=^$ -cpu 1,2,4,8
go test -v -test.bench=BenchmarkIp -benchmem -run=^$
*/

func TestTcpConn_NewTcpConnServer(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	//Sever
	go func() {
		defer wg.Done()
		//o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(ClientOrServer(SERVER))
		assert.NoError(t, err)
		DeLog.Infof(INFOPREFIX+"server,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	wg.Add(1)
	time.Sleep(3 * time.Second)
	//Client
	go func() {
		defer wg.Done()
		//o := newOptions()
		con, err := newTcpConn()
		assert.NoError(t, err)
		DeLog.Infof(INFOPREFIX+"client,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	wg.Wait()
}
func TestTcpConn_NewTcpConnClient(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	//Client
	go func() {
		defer wg.Done()
		//o := newOptions()
		con, err := newTcpConn()
		assert.NoError(t, err)
		DeLog.Infof(INFOPREFIX+"client,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	//Sever
	wg.Add(1)
	time.Sleep(3 * time.Second)
	go func() {
		defer wg.Done()
		//o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(ClientOrServer(SERVER))
		assert.NoError(t, err)
		DeLog.Infof(INFOPREFIX+"server,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	wg.Wait()
}
func TestTcpConn_TimeOutServer(t *testing.T) {
	//o := newOptions(ClientOrServer(SERVER), TimeOut(3))
	con, err := newTcpConn(ClientOrServer(SERVER), TimeOut(3))
	assert.Error(t, err, "error ret:%v", err)
	DeLog.Infof(INFOPREFIX+"server,con:%+v\n", con)
}
func TestTcpConn_TimeOutClient(t *testing.T) {
	//o := newOptions(ClientOrServer(CLIENT), TimeOut(3))
	con, err := newTcpConn(ClientOrServer(CLIENT), TimeOut(3))
	assert.Error(t, err, "error ret:%v", err)
	DeLog.Infof(INFOPREFIX+"client,con:%+v\n", con)
}
func TestTcpConn_Close(t *testing.T) {
	go func() {
		//o := newOptions(ClientOrServer(CLIENT))
		con, err := newTcpConn(ClientOrServer(CLIENT))
		assert.NoError(t, err)
		DeLog.Infof(INFOPREFIX+"client,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	time.Sleep(time.Second)
	//o := newOptions(ClientOrServer(SERVER))
	con, err := newTcpConn(ClientOrServer(SERVER))
	assert.NoError(t, err)
	DeLog.Infof(INFOPREFIX+"server,con:%+v\n", con)
	err = con.Close()
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)
	go func() {
		time.Sleep(time.Second)
		//o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(ClientOrServer(SERVER))
		assert.NoError(t, err)
		DeLog.Infof(INFOPREFIX+"server,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	//o1 := newOptions(ClientOrServer(CLIENT))
	con1, err1 := newTcpConn(ClientOrServer(CLIENT))
	assert.NoError(t, err1)
	DeLog.Infof(INFOPREFIX+"client,con:%+v\n", con)
	err1 = con1.Close()
	assert.NoError(t, err1)
}
func genRandData(seed int, datasetNum int, dataSrcLen int) ([]byte, []int) {
	dataSrc := make([]byte, dataSrcLen)
	//b := []byte("s")[0]
	rand.Seed(int64(seed))
	for i := 0; i < len(dataSrc); i++ {
		dataSrc[i] = uint8(rand.Uint32() % 256)
	}
	offList := make([]int, datasetNum)
	for j := 0; j < datasetNum; j++ {
		off := int(rand.Uint32()) % len(dataSrc)
		if off != 0 {
			offList[j] = off
		} else {
			offList[j] = 7
		}
	}
	return dataSrc, offList
}

func TestTcpConn_ClientRecvData(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	test_num := 10000
	srcData, offList := genRandData(111, test_num, 4096)
	DeLog.Infof(INFOPREFIX+"offList:%+v\n", offList[:100])
	go func() {
		defer wg.Done()
		//o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(ClientOrServer(SERVER))
		assert.NoError(t, err)
		defer func() {
			err := con.Close()
			assert.NoError(t, err)
		}()
		for i := 0; i < test_num; i++ {
			perData := srcData[:offList[i]]
			n, err := con.SendData("", perData)
			assert.NoError(t, err)
			assert.Equal(t, n, len(perData))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		//o := newOptions(ClientOrServer(CLIENT))
		con, err := newTcpConn(ClientOrServer(CLIENT))
		assert.NoError(t, err)
		defer func() {
			err := con.Close()
			assert.NoError(t, err)
		}()
		for i := 0; i < test_num; i++ {
			dataR, err := con.RecvData("")
			assert.NoError(t, err, "***err:%+v", err)
			actualData := srcData[:offList[i]]
			assert.Equal(t, len(actualData), len(dataR), "len(actualData):%v, len(dataR):%v", len(actualData), len(dataR))
			assert.Equal(t, actualData, dataR)
			//fmt.Printf("i:%+v,client:%+v\n", i, dataR[:7])
		}
	}()
	wg.Wait()
}

func TestTcpConn_ClientSendData(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	test_num := 10000
	srcData, offList := genRandData(111, test_num, 4096)
	DeLog.Infof(INFOPREFIX+"offList:%+v\n", offList[:100])
	go func() {
		defer wg.Done()
		//o := newOptions(ClientOrServer(CLIENT))
		con, err := newTcpConn(ClientOrServer(CLIENT))
		assert.NoError(t, err)
		defer func() {
			err := con.Close()
			assert.NoError(t, err)
		}()
		for i := 0; i < test_num; i++ {
			perData := srcData[:offList[i]]
			n, err := con.SendData("", perData)
			assert.NoError(t, err)
			assert.Equal(t, n, len(perData))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		//o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(ClientOrServer(SERVER))
		assert.NoError(t, err)
		defer func() {
			err := con.Close()
			assert.NoError(t, err)
		}()
		for i := 0; i < test_num; i++ {
			dataR, err := con.RecvData("")
			assert.NoError(t, err, "***err:%+v", err)
			actualData := srcData[:offList[i]]
			assert.Equal(t, len(actualData), len(dataR), "len(actualData):%v, len(dataR):%v", len(actualData), len(dataR))
			assert.Equal(t, actualData, dataR)
			//fmt.Printf("i:%+v,client:%+v\n", i, dataR[:7])
		}
	}()
	wg.Wait()
}

func TestTcpConn_RecvDataExitServer(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		//o := newOptions(ClientOrServer(CLIENT))
		con, err := newTcpConn(ClientOrServer(CLIENT))
		assert.NoError(t, err)
		tick := time.NewTicker(5 * time.Second)
		defer func() {
			err := con.Close()
			assert.NoError(t, err)
			tick.Stop()
		}()
		<-tick.C
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		//o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(ClientOrServer(SERVER))
		assert.NoError(t, err)
		defer func() {
			err := con.Close()
			assert.NoError(t, err)
		}()
		DeLog.Infof(INFOPREFIX + "server: RecvData start\n")
		dataR, err := con.RecvData("")
		//assert.NoError(t, err, "***err:%+v", err)
		DeLog.Infof(INFOPREFIX+"server: RecvData exit,dataR:%+v,err:%+v\n", dataR, err)
	}()
	wg.Wait()
}

func TestTcpConn_RecvDataExitClient(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		//o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(ClientOrServer(SERVER))
		assert.NoError(t, err)
		tick := time.NewTicker(5 * time.Second)
		defer func() {
			err := con.Close()
			assert.NoError(t, err)
			tick.Stop()
		}()
		<-tick.C
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		//o := newOptions(ClientOrServer(CLIENT))
		con, err := newTcpConn(ClientOrServer(CLIENT))
		assert.NoError(t, err)
		defer func() {
			err := con.Close()
			assert.NoError(t, err)
		}()
		DeLog.Infof(INFOPREFIX + "client: RecvData start\n")
		dataR, err := con.RecvData("")
		//assert.NoError(t, err, "***err:%+v", err)
		DeLog.Infof(INFOPREFIX+"client: RecvData exit,dataR:%+v,err:%+v\n", dataR, err)
	}()
	wg.Wait()
}

func BenchmarkConcatStringByAdd(b *testing.B) {
	elems := []string{"1", "2", "3", "4", "5"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret := ""
		for _, elem := range elems {
			ret += elem
		}
	}
	b.StopTimer()
}
