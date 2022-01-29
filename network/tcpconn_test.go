package network

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

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
		o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(o)
		assert.NoError(t, err)
		fmt.Printf("***server,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	wg.Add(1)
	time.Sleep(3 * time.Second)
	//Client
	go func() {
		defer wg.Done()
		o := newOptions()
		con, err := newTcpConn(o)
		assert.NoError(t, err)
		fmt.Printf("***client,con:%+v\n", con)
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
		o := newOptions()
		con, err := newTcpConn(o)
		assert.NoError(t, err)
		fmt.Printf("***client,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	//Sever
	wg.Add(1)
	time.Sleep(3 * time.Second)
	go func() {
		defer wg.Done()
		o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(o)
		assert.NoError(t, err)
		fmt.Printf("***server,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	wg.Wait()
}
func TestTcpConn_TimeOutServer(t *testing.T) {
	o := newOptions(ClientOrServer(SERVER), TimeOut(3))
	con, err := newTcpConn(o)
	assert.Error(t, err, "error ret:%v", err)
	fmt.Printf("***server,con:%+v\n", con)
}
func TestTcpConn_TimeOutClient(t *testing.T) {
	o := newOptions(ClientOrServer(CLIENT), TimeOut(3))
	con, err := newTcpConn(o)
	assert.Error(t, err, "error ret:%v", err)
	fmt.Printf("***client,con:%+v\n", con)
}
func TestTcpConn_Close(t *testing.T) {
	go func() {
		o := newOptions(ClientOrServer(CLIENT))
		con, err := newTcpConn(o)
		assert.NoError(t, err)
		fmt.Printf("***client,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	time.Sleep(time.Second)
	o := newOptions(ClientOrServer(SERVER))
	con, err := newTcpConn(o)
	assert.NoError(t, err)
	fmt.Printf("***server,con:%+v\n", con)
	err = con.Close()
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)
	go func() {
		time.Sleep(time.Second)
		o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(o)
		assert.NoError(t, err)
		fmt.Printf("***server,con:%+v\n", con)
		err = con.Close()
		assert.NoError(t, err)
	}()
	o1 := newOptions(ClientOrServer(CLIENT))
	con1, err1 := newTcpConn(o1)
	assert.NoError(t, err1)
	fmt.Printf("***client,con:%+v\n", con1)
	err1 = con1.Close()
	assert.NoError(t, err1)
}
func genRandData(seed int, datasetNum int) ([]byte, []int) {
	dataSrc := make([]byte, 4096)
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
	srcData, offList := genRandData(111, test_num)
	fmt.Printf("offList:%+v\n", offList[:100])
	go func() {
		defer wg.Done()
		o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(o)
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
		o := newOptions(ClientOrServer(CLIENT))
		con, err := newTcpConn(o)
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
	srcData, offList := genRandData(111, test_num)
	fmt.Printf("offList:%+v\n", offList[:100])
	go func() {
		defer wg.Done()
		o := newOptions(ClientOrServer(CLIENT))
		con, err := newTcpConn(o)
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
		o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(o)
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
		o := newOptions(ClientOrServer(CLIENT))
		con, err := newTcpConn(o)
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
		o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(o)
		assert.NoError(t, err)
		defer func() {
			err := con.Close()
			assert.NoError(t, err)
		}()
		fmt.Printf("***server: RecvData start\n")
		dataR, err := con.RecvData("")
		//assert.NoError(t, err, "***err:%+v", err)
		fmt.Printf("***server: RecvData exit,dataR:%+v,err:%+v\n", dataR, err)
	}()
	wg.Wait()
}

func TestTcpConn_RecvDataExitClient(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		o := newOptions(ClientOrServer(SERVER))
		con, err := newTcpConn(o)
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
		o := newOptions(ClientOrServer(CLIENT))
		con, err := newTcpConn(o)
		assert.NoError(t, err)
		defer func() {
			err := con.Close()
			assert.NoError(t, err)
		}()
		fmt.Printf("***client: RecvData start\n")
		dataR, err := con.RecvData("")
		//assert.NoError(t, err, "***err:%+v", err)
		fmt.Printf("***client: RecvData exit,dataR:%+v,err:%+v\n", dataR, err)
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
