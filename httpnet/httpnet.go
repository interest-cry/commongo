package httpnet

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var allAlgoId = []string{"algo_0", "algo_1", "algo_2"}

// var xx = map[string]int{}

type register struct {
	names       []string
	nameToIndex map[string]int
	mapList     []*sync.Map
}

func newRegister() *register {
	r := register{nameToIndex: make(map[string]int)}
	for i := 0; i < len(allAlgoId); i++ {
		r.names = append(r.names, allAlgoId[i])
		r.nameToIndex[allAlgoId[i]] = i
		r.mapList = append(r.mapList, &sync.Map{})
	}
	return &r
}
func (r *register) getMap(algoName string) *sync.Map {
	ind, ok := r.nameToIndex[algoName]
	if !ok {
		panic(errors.New("error:no algo name"))
	}
	return r.mapList[ind]
}

type DataInfo struct {
	Data interface{}
}

// type algoEleMap struct{ sidMap sync.Map }
type Net struct {
	engine  *gin.Engine
	algoReg *register
}
type Connect struct {
	algoId     string
	taskId     string
	sid        string
	uid        string
	remoteAddr string
	n          *Net
	m          *sync.Map
	client     *http.Client
}

func NewConnect(algoId string, taskId string, sid string, n *Net,
	remoteIp string, remotePort int) *Connect {
	uid := taskId + "_" + sid
	return &Connect{algoId: algoId, taskId: taskId, sid: sid,
		uid: uid, remoteAddr: "http://" + remoteIp + ":" + strconv.Itoa(remotePort) + "/send",
		n: n, m: n.algoReg.getMap(algoId),
		client: &http.Client{Transport: http.DefaultTransport}}
}
func (c *Connect) Close() {
	c.m.Delete(c.uid)
}
func (c *Connect) Recv() *NetInfo {
	for {
		v, ok := c.m.Load(c.uid)
		if ok {
			ch := v.(chan *NetInfo)
			r := <-ch
			return r
		}
		// time.Sleep(time.Second)
	}
}
func (c *Connect) Send(value string) error {
	netInfo := NetInfo{AlgoId: c.algoId, TaskId: c.taskId, Sid: c.sid, Value: value}
	dataJson, _ := json.Marshal(netInfo)
	// fmt.Printf("c.remoteAddr:%v\n", c.remoteAddr)
	rsp, err := c.client.Post(c.remoteAddr, "application/json", bytes.NewReader(dataJson))
	if err != nil {
		return err
	}
	rsp.Body.Close()
	return nil
}
func NewNet(ip string, port int) *Net {
	var n = Net{engine: gin.New(), algoReg: newRegister()}
	n.engine.POST("/send", n.pro_send)
	go func() {
		if err := n.engine.Run(ip + ":" + strconv.Itoa(port)); err != nil {
			panic(err)
		}
	}()
	return &n
}

type NetInfo struct {
	AlgoId string
	TaskId string
	Sid    string
	Key    string
	Value  string
}

func (n *Net) pro_send(c *gin.Context) {
	var netInfo NetInfo
	c.ShouldBindJSON(&netInfo)
	// fmt.Printf("pro_send>>>netInfo:%+v\n", netInfo)
	uid := netInfo.TaskId + "_" + netInfo.Sid
	m := n.algoReg.getMap(netInfo.AlgoId)
	ch1, ok1 := m.Load(uid)
	// fmt.Printf("pro_send>>>ch1, ok1:%+v,%v\n", ch1, ok1)
	if ok1 {
		ch := ch1.(chan *NetInfo)
		ch <- &netInfo
	} else {
		ch := make(chan *NetInfo)
		m.Store(uid, ch)
		ch <- &netInfo
	}
}
