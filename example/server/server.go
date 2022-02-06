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
	NetworkType string `json:"network_type"`
	Role        string `json:"role"`
	SendUrl     string `json:"send_url"`
	Uid         string `json:"uid"`
	LocalNid    string `json:"local_nid"`
	RemoteNid   string `json:"remote_nid"`
	//HostRemotes []string `json:"host_remotes"`
	//Host string
}
type Server struct {
	HttpBigcache *network.HttpBigCache
	Eventbus     *network.EventBus
	Addr         string
	NetworkType  string
}

func NewServer(addr, networkType string) (*Server, error) {
	//conf := bigcache.DefaultConfig(1800 * time.Second)
	//conf.CleanWindow = time.Millisecond * 500
	//fmt.Printf("bigcache config:%+v\n", conf)
	//bigC, err := bigcache.NewBigCache(conf)
	//httpBigCache := network.NewHttpBigCache(1800)
	httpBigCache := network.DefaultHttpBigCache
	//Eventbus := network.DefaultEventBus
	Eventbus := network.NewEventBus(60)
	return &Server{
		Addr:         addr,
		HttpBigcache: httpBigCache,
		Eventbus:     Eventbus,
		NetworkType:  networkType}, nil
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

func (s *Server) StartTask(c *gin.Context) {
	req := Req{}
	err := c.BindJSON(&req)
	if err != nil {
		network.DeLog.Infof(network.INFOPREFIX+"SaveData BindJson error:%v\n", err)
		return
	}
	network.DeLog.Infof("===req:%+v", req)
	datasetNum := 1000 * 100
	//datasetNum = 10
	dataSrcLen := 10240000
	switch req.Role {
	case GUEST:
		msgHandle, err := network.NewMessager(
			network.NetworkMap[req.NetworkType],
			network.BigCache(s.HttpBigcache),
			network.SendUrl(req.SendUrl),
			network.Uid(req.Uid),
			network.EventBusSet(s.Eventbus),
			network.LocalNid(req.LocalNid),
			network.RemoteNid(req.RemoteNid))
		if err != nil {
			c.JSON(200, &Rsp{err.Error(), -10})
			return
		}
		defer func() {
			msgHandle.Close()
		}()
		//offsetList:=[]int
		srcData, _ := GenRandDataDebug(111, datasetNum, dataSrcLen)
		for i := 0; i < datasetNum; i++ {
			//先发送
			send_key := req.LocalNid + "_" + req.Uid + "_" + strconv.Itoa(i)
			//fmt.Printf("+++++++++++>>send_key:%v\n", send_key)
			srcData = []byte("msg_" + send_key)
			_, err := msgHandle.SendData(send_key, srcData)
			if err != nil {
				network.DeLog.Infof(network.INFOPREFIX+"guest send data error:%v\n", err)
				c.JSON(200, gin.H{"msg": err.Error()})
				return
			}
			//再接收
			key := req.RemoteNid + "_" + req.Uid + "_" + strconv.Itoa(i)
			data, err := msgHandle.RecvData(key)
			if err != nil {
				network.DeLog.Infof(network.INFOPREFIX+"guest recv data error:%v\n", err)
				c.JSON(200, gin.H{"msg": err.Error()})
				return
			}
			network.DeLog.Infof("[%v]RecvData,key:%v,data:%v", req.LocalNid, key, string(data))
		}
	case HOST:
		msgHandle, err := network.NewMessager(
			network.NetworkMap[req.NetworkType],
			network.BigCache(s.HttpBigcache),
			network.SendUrl(req.SendUrl),
			network.Uid(req.Uid),
			network.EventBusSet(s.Eventbus),
			network.LocalNid(req.LocalNid),
			network.RemoteNid(req.RemoteNid))
		if err != nil {
			c.JSON(200, &Rsp{err.Error(), -10})
			return
		}
		defer func() {
			msgHandle.Close()
		}()
		//ret, err := msgHandle.RecvData("key")
		for i := 0; i < datasetNum; i++ {
			//先接收
			key := req.RemoteNid + "_" + req.Uid + "_" + strconv.Itoa(i)
			//fmt.Printf("+++++++++++>>key:%v\n", key)
			data, err := msgHandle.RecvData(key)
			if err != nil {
				network.DeLog.Infof(network.INFOPREFIX+"host recv data error:%v\n", err)
				c.JSON(200, gin.H{"msg": err.Error()})
				return
			}
			//network.DeLog.Infof(network.INFOPREFIX+"i:%v,ret len:%v,err:%v\n", i, len(ret), err)
			network.DeLog.Infof("[%v]RecvData,key:%v,data:%v", req.LocalNid, key, string(data))
			//再发送
			send_key := req.LocalNid + "_" + req.Uid + "_" + strconv.Itoa(i)
			_, err = msgHandle.SendData(send_key, []byte("msg_"+send_key))
			if err != nil {
				network.DeLog.Infof(network.INFOPREFIX+"host send data error:%v\n", err)
				c.JSON(200, gin.H{"msg": err.Error()})
				return
			}
		}
	}
	c.JSON(200, gin.H{"msg": "ok"})
	return
}
func (s *Server) Route() {
	router := gin.New()
	//router := gin.Default()
	v1Grp := router.Group("/v1")
	v1Grp.POST("/send", s.HttpBigcache.HttpBigCacheHandlerFunc)
	//v1Grp.POST("/send", s.Eventbus.EventBusHandlerFunc)
	v1Grp.POST("/algo/start", s.StartTask)
	if err := router.Run(s.Addr); err != nil {
		network.DeLog.Infof(network.INFOPREFIX + "server run error")
	}
}
