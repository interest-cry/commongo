package network

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strconv"
	"testing"
)

func mockHttpConn(sendUrl string, hbigC *HttpBigCache) *HttpConn {
	//bigC, err := bigcache.NewBigCache(bigcache.DefaultConfig(3600 * time.Second))
	o := newOptions(BigCache(hbigC), SendUrl(sendUrl))
	httpcache, err := newHttpConn(o)
	if err != nil {
		panic(err)
	}
	return httpcache
}

type mockServer struct {
	Addr  string
	HBigC *HttpBigCache
}

func newMockServer(addr string) (*mockServer, error) {
	//bigC, err := bigcache.NewBigCache(bigcache.DefaultConfig(30 * time.Second))
	//hBigC := NewHttpBigCache(30)
	hBigC := DefaultHttpBigCache
	s := new(mockServer)
	s.HBigC = hBigC
	s.Addr = addr
	return s, nil
}
func (s *mockServer) saveData(c *gin.Context) {
	//ser.SaveData(c)
	var req HttpConnRequest
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Printf("SaveData BindJson error:%v\n", err)
		return
	}
	err = s.HBigC.bigC.Set(req.Key, req.Data)
	if err != nil {
		fmt.Printf("SaveData set val error:%v", err)
		return
	}
	fmt.Printf("save data ok\n")
	//message.Log.Infof("===>>本地缓存 set data ok")
	//todo:不需要响应返回
	//c.JSON(200, gin.H{
	//	"msg": "ok",
	//})

}
func mockServerRun(relativePath string) (*gin.Engine, *mockServer) {
	r := gin.Default()
	ms, err := newMockServer(":8900")
	if err != nil {
		panic(err)
	}
	r.POST(relativePath, ms.saveData)
	return r, ms
}
func TestHttpCache_SendData(t *testing.T) {
	relativePath := "/v1/send"
	//mockSer, err := newMockServer(relativePath)
	rout, ms := mockServerRun(relativePath)
	service := httptest.NewServer(rout)
	defer func() { service.Close() }()
	httpCache := mockHttpConn(service.URL+relativePath, ms.HBigC)
	_, err := httpCache.SendData("kk", []byte("yyyy"))
	assert.NoError(t, err)
	ret, err := httpCache.RecvData("kk")
	assert.NoError(t, err)
	fmt.Printf("ret:%v\n", string(ret))
	//send data
	//wg:=sync.WaitGroup{}
	seed := 111
	datasetNum := 9770
	dataSrc, _ := genRandData(seed, datasetNum, 102400)
	go func() {
		for i := 0; i < datasetNum; i++ {
			key := strconv.Itoa(i)
			_, err := httpCache.SendData(key, dataSrc)
			assert.NoError(t, err)
		}
	}()
	for i := 0; i < datasetNum; i++ {
		key := strconv.Itoa(i)
		ret, err := httpCache.RecvData(key)
		assert.NoError(t, err)
		assert.Equal(t, ret, dataSrc)
	}
	fmt.Printf("==========================end\n")
	//time.Sleep(3600 * time.Second)
}
