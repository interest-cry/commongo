package server

import (
	"bytes"
	"encoding/json"
	"github.com/Zhenhanyijiu/commongo/network"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

const (
	GUEST = "GUEST"
	HOST  = "HOST"
)

type Rsp struct {
	ErrMsg  string `json:"err_msg"`
	ErrCode int    `json:"err_code"`
}
type Req struct {
	Role    string `json:"role"`
	SendUrl string `json:"send_url"`
	//HostRemotes []string `json:"host_remotes"`
	//Host string
}
type Server struct {
	Addr  string
	HBigC *network.HttpBigCache
}

func StartTestSdk(req *Req, reqUrl string) {
	dataJson, _ := json.Marshal(req)
	httpClient := http.Client{Transport: http.DefaultTransport}
	rsp, err := httpClient.Post(reqUrl, "application/json", bytes.NewReader(dataJson))
	if err != nil {
		network.DeLog.Infof(network.INFOPREFIX+"post error:%+v\n", err)
		return
	}
	defer func() {
		rsp.Body.Close()
	}()
	ret, err := ioutil.ReadAll(rsp.Body)
	network.DeLog.Infof(network.INFOPREFIX+"rsp:%v,err:%v\n", string(ret), err)
	return
}
func GenRandDataDebug(seed int, datasetNum int, dataSrcLen int) ([]byte, []int) {
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

func NewServer(addr string) (*Server, error) {
	//conf := bigcache.DefaultConfig(1800 * time.Second)
	//conf.CleanWindow = time.Millisecond * 500
	//fmt.Printf("bigcache config:%+v\n", conf)
	//bigC, err := bigcache.NewBigCache(conf)
	//httpBigCache := network.NewHttpBigCache(1800)
	httpBigCache := network.DefaultHttpBigCache
	s := new(Server)
	s.HBigC = httpBigCache
	s.Addr = addr
	return s, nil
}

func (s *Server) StartTask(c *gin.Context) {
	req := Req{}
	err := c.BindJSON(&req)
	if err != nil {
		network.DeLog.Infof(network.INFOPREFIX+"SaveData BindJson error:%v\n", err)
		return
	}
	msgHandle, err := network.NewMessager(network.HTTPCACHE,
		network.BigCache(s.HBigC),
		network.SendUrl(req.SendUrl))
	if err != nil {
		c.JSON(200, &Rsp{err.Error(), -10})
		return
	}
	defer func() {
		msgHandle.Close()
	}()

	datasetNum := 977 * 5
	//datasetNum = 10
	dataSrcLen := 102400
	switch req.Role {
	case GUEST:
		//offsetList:=[]int
		srcData, _ := GenRandDataDebug(111, datasetNum, dataSrcLen)
		for i := 0; i < datasetNum; i++ {
			key := "key_" + strconv.Itoa(i)
			srcData = []byte(key)
			_, err := msgHandle.SendData(key, srcData)
			if err != nil {
				network.DeLog.Infof(network.INFOPREFIX+"guest send data error:%v\n", err)
				c.JSON(200, gin.H{"msg": err.Error()})
				return
			}
			network.DeLog.Infof("SendData,i:%v,key:%v,data:%v", i, key, string(srcData))
		}
	case HOST:
		//ret, err := msgHandle.RecvData("key")
		for i := 0; i < datasetNum; i++ {
			key := "key_" + strconv.Itoa(i)
			ret, err := msgHandle.RecvData(key)
			if err != nil {
				network.DeLog.Infof(network.INFOPREFIX+"host recv data error:%v\n", err)
				c.JSON(200, gin.H{"msg": err.Error()})
				return
			}
			//network.DeLog.Infof(network.INFOPREFIX+"i:%v,ret len:%v,err:%v\n", i, len(ret), err)
			network.DeLog.Infof("RecvData,i:%v,key:%v,data:%v", i, key, string(ret))
		}
	}
	c.JSON(200, gin.H{"msg": "ok"})
	return
}
func (s *Server) Route() {
	router := gin.New()
	v1Grp := router.Group("/v1")
	v1Grp.POST("/send", s.HBigC.BigCacheHandlerFunc)
	v1Grp.POST("/algo/start", s.StartTask)
	if err := router.Run(s.Addr); err != nil {
		network.DeLog.Infof(network.INFOPREFIX + "server run error")
	}
}
