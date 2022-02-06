package network

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

//go test -v -bench=BenchmarkIp -benchmem -run=^$

type mockHttpServer struct {
	httpBigCache *HttpBigCache
	engine       *gin.Engine
}

func newMockHttpServer(sec int) *mockHttpServer {
	httpBigCache := NewHttpBigCache(sec)
	//httpBigCache := DefaultHttpBigCache
	return &mockHttpServer{
		httpBigCache: httpBigCache,
		engine:       gin.New(),
	}
}

func (h *mockHttpServer) addPath(relativePath string) {
	h.engine.POST(relativePath, h.httpBigCache.HttpBigCacheHandlerFunc)
}

//func (s *mockServer) saveData(c *gin.Context) {
//	//ser.SaveData(c)
//	var req HttpConnRequest
//	err := c.BindJSON(&req)
//	if err != nil {
//		DeLog.Infof(INFOPREFIX+"SaveData BindJson error:%v", err)
//		return
//	}
//	err = s.HBigC.bigC.Set(req.Key, req.Data)
//	if err != nil {
//		DeLog.Infof(INFOPREFIX+"SaveData set val error:%v", err)
//		return
//	}
//	DeLog.Infof(INFOPREFIX + "save data ok")
//	//message.Log.Infof("===>>本地缓存 set data ok")
//	//todo:不需要响应返回
//	//c.JSON(200, gin.H{
//	//	"msg": "ok",
//	//})
//
//}
func TestNewHttpBigCache(t *testing.T) {
	NewHttpBigCache(10)
}
func TestCacheConn_Close(t *testing.T) {

}
func TestHttpBigCache_HttpBigCacheHandlerFunc(t *testing.T) {
	//启动第一个服务
	hServer1 := newMockHttpServer(60)
	path1 := "/v1/send"
	nid1 := "nid1"
	hServer1.addPath(path1)
	service1 := httptest.NewServer(hServer1.engine)
	defer func() { service1.Close() }()
	//启动第二个服务
	hServer2 := newMockHttpServer(60)
	path2 := "/v1/send"
	nid2 := "nid2"
	hServer2.addPath(path2)
	service2 := httptest.NewServer(hServer2.engine)
	defer func() { service2.Close() }()
	wg := sync.WaitGroup{}
	wg.Add(2)
	datasetNum := 100

	//服务1
	go func() {
		defer func() { wg.Done() }()
		//创建一个连接
		httpConn2, err := newHttpConn(
			SendUrl(service2.URL+path2),
			BigCache(hServer1.httpBigCache))
		assert.NoError(t, err)
		defer func() { httpConn2.Close() }()
		for i := 0; i < datasetNum; i++ {
			//先接收
			key := nid2 + "_" + strconv.Itoa(i)
			data, err := httpConn2.RecvData(key)
			assert.NoError(t, err)
			DeLog.Infof("[%v],key:%v,data:%v,RecvData", nid1, key, string(data))
			//再发送
			send_key := nid1 + "_" + strconv.Itoa(i)
			httpConn2.SendData(send_key, []byte("msg_"+send_key))
		}
	}()
	//服务2
	go func() {
		defer func() { wg.Done() }()
		//创建一个连接
		httpConn1, err := newHttpConn(
			SendUrl(service1.URL+path1),
			BigCache(hServer2.httpBigCache))
		assert.NoError(t, err)
		defer func() { httpConn1.Close() }()
		for i := 0; i < datasetNum; i++ {
			//先发送
			send_key := nid2 + "_" + strconv.Itoa(i)
			httpConn1.SendData(send_key, []byte("msg_"+send_key))
			//再接收
			key := nid1 + "_" + strconv.Itoa(i)
			data, err := httpConn1.RecvData(key)
			assert.NoError(t, err)
			DeLog.Infof("[%v],key:%v,data:%v,RecvData", nid2, key, string(data))
		}
	}()
	wg.Wait()
}

func TestCacheConn_RecvData(t *testing.T) {

}
func TestCacheConn_SendData(t *testing.T) {
	relativePath := "/v1/send"
	//mockSer, err := newMockServer(relativePath)
	httpServer := newMockHttpServer(60)
	httpServer.addPath(relativePath)
	//启动一个服务
	service := httptest.NewServer(httpServer.engine)
	defer func() { service.Close() }()
	httpConn, err := newHttpConn(SendUrl(service.URL+relativePath),
		BigCache(httpServer.httpBigCache))
	assert.NoError(t, err)
	defer func() { httpConn.Close() }()
	_, err = httpConn.SendData("kk", []byte("yyyy"))
	assert.NoError(t, err)
	ret, err := httpConn.RecvData("kk")
	assert.NoError(t, err)
	DeLog.Infof(INFOPREFIX+"ret:%v", string(ret))
	//send data
	//wg:=sync.WaitGroup{}
	seed := 111
	datasetNum := 977
	dataSrc, _ := genRandData(seed, datasetNum, 102400)
	go func() {
		for i := 0; i < datasetNum; i++ {
			key := strconv.Itoa(i)
			_, err := httpConn.SendData(key, dataSrc)
			assert.NoError(t, err)
		}
	}()
	for i := 0; i < datasetNum; i++ {
		key := strconv.Itoa(i)
		ret, err := httpConn.RecvData(key)
		assert.NoError(t, err)
		assert.Equal(t, ret, dataSrc)
	}
	DeLog.Infof(INFOPREFIX + "===end")
	//time.Sleep(3600 * time.Second)
}

func BenchmarkHttpBigCache_HttpBigCacheHandlerFunc(b *testing.B) {
	//启动第一个服务
	hServer1 := newMockHttpServer(60)
	path1 := "/v1/send"
	nid1 := "nid1"
	hServer1.addPath(path1)
	service1 := httptest.NewServer(hServer1.engine)
	defer func() { service1.Close() }()
	//启动第二个服务
	hServer2 := newMockHttpServer(60)
	path2 := "/v1/send"
	nid2 := "nid2"
	hServer2.addPath(path2)
	service2 := httptest.NewServer(hServer2.engine)
	defer func() { service2.Close() }()
	wg := sync.WaitGroup{}
	wg.Add(2)
	//datasetNum := 100
	//t := &testing.T{}
	b.ResetTimer()
	//b.N = 10
	//服务1
	go func() {
		defer func() { wg.Done() }()
		//创建一个连接
		httpConn2, _ := newHttpConn(
			SendUrl(service2.URL+path2),
			BigCache(hServer1.httpBigCache))
		//assert.NoError(t, err)
		defer func() { httpConn2.Close() }()
		for i := 0; i < b.N; i++ {
			//先接收
			key := nid2 + "_" + strconv.Itoa(i)
			data, _ := httpConn2.RecvData(key)
			//assert.NoError(t, err)
			//DeLog.Infof("[%v],key:%v,data:%v,RecvData", nid1, key, string(data))
			bytes.Equal(data, []byte("msg_"+key))
			//再发送
			send_key := nid1 + "_" + strconv.Itoa(i)
			httpConn2.SendData(send_key, []byte("msg_"+send_key))
		}
	}()
	//服务2
	go func() {
		defer func() { wg.Done() }()
		//创建一个连接
		httpConn1, _ := newHttpConn(
			SendUrl(service1.URL+path1),
			BigCache(hServer2.httpBigCache))
		//assert.NoError(t, err)
		defer func() { httpConn1.Close() }()
		for i := 0; i < b.N; i++ {
			//先发送
			send_key := nid2 + "_" + strconv.Itoa(i)
			httpConn1.SendData(send_key, []byte("msg_"+send_key))
			//再接收
			key := nid1 + "_" + strconv.Itoa(i)
			data, _ := httpConn1.RecvData(key)
			//assert.NoError(t, err)
			//DeLog.Infof("[%v],key:%v,data:%v,RecvData", nid2, key, string(data))
			bytes.Equal(data, []byte("msg_"+key))
		}
	}()
	wg.Wait()
}
