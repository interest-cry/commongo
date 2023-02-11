package httpnetgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	//SERVER           = "server"
	//CLIENT           = "client"
	REMOTEIP         = "ip"
	REMOTEPORT       = "port"
	APPNAME          = "appname"
	SESSID           = "sessid"
	UniqueServerName = "sername"
	INITCONNPATH     = "/initconn"
	SENDPATH         = "/send"
)

var (
	LocalDebug = false
)

type webApps struct {
	webAppNames []string
	nameToIndex map[string]int
	mapList     []*sync.Map
}

func newWebApps(webAppNames []string) *webApps {
	w := webApps{nameToIndex: make(map[string]int)}
	for i := 0; i < len(webAppNames); i++ {
		w.webAppNames = webAppNames
		w.nameToIndex[webAppNames[i]] = i
		w.mapList = append(w.mapList, &sync.Map{})
	}
	return &w
}
func (w *webApps) isAppNameIn(appname string) (index int, ok bool) {
	ind, ok := w.nameToIndex[appname]
	return ind, ok
}
func (w *webApps) getMap(index int) *sync.Map {
	return w.mapList[index]
}

type NetServer struct {
	uniqueServerName string
	appNames         *webApps
	ip               string
	port             string
	engine           *gin.Engine
	sessId           uint32
	mutexWR          sync.Mutex
}

func NewNetServer(ipAnPort string, uniqueServerName string, webAppNames []string) *NetServer {
	ret := strings.Split(ipAnPort, ":")
	if len(ret) != 2 {
		panic("error:ip and port is wrong.please use ip:port(127.0.0.1:18888)")
	}
	var n = NetServer{
		uniqueServerName: uniqueServerName,
		appNames:         newWebApps(webAppNames),
		ip:               ret[0],
		port:             ret[1],
		engine:           gin.New(),
		sessId:           0}
	if len(webAppNames) == 0 || uniqueServerName == "" {
		panic("error:web app name and server name must not null")
	}
	n.engine.GET(INITCONNPATH, n.initconncall)
	n.engine.POST(SENDPATH, n.sendcall)
	go func() {
		addr := ":" + n.port
		if LocalDebug {
			addr = n.ip + ":" + n.port
		}
		if err := n.engine.Run(addr); err != nil {
			panic(err)
		}
	}()
	return &n
}

func (n *NetServer) initconncall(c *gin.Context) {
	n.mutexWR.Lock()
	defer func() { n.mutexWR.Unlock() }()
	sid := atomic.AddUint32(&n.sessId, 1)
	appname := c.GetHeader(APPNAME)
	remoteIp := c.GetHeader(REMOTEIP)
	remoteport := c.GetHeader(REMOTEPORT)
	ind, ok := n.appNames.isAppNameIn(appname)
	if !ok {
		c.JSON(200, &NetRspInfo{ErrCode: 1, ErrMsg: "initconn failed,appname not exist!"})
		return
	}
	if remoteIp == "" || remoteport == "" {
		c.JSON(200, &NetRspInfo{ErrCode: 1, ErrMsg: "initconn failed,ip or port wrong!"})
		return
	}
	m := n.appNames.getMap(ind)
	ch := make(chan *NetReqInfo)
	sidString := strconv.Itoa(int(sid))
	m.Store(sidString, &MemoryInfo{ch: ch, remoteIp: remoteIp, remoteport: remoteport})
	c.Writer.Header().Set(SESSID, sidString)
	c.Writer.Header().Set(UniqueServerName, n.uniqueServerName)
	c.JSON(200, &NetRspInfo{ErrCode: 0, ErrMsg: "ok"})
}
func (n *NetServer) sendcall(c *gin.Context) {
	var netInfo NetReqInfo
	c.ShouldBindJSON(&netInfo)
	sid := c.GetHeader(SESSID)
	sername := c.GetHeader(UniqueServerName)
	appname := c.GetHeader(APPNAME)
	if sername != n.uniqueServerName {
		c.JSON(200, &NetRspInfo{ErrCode: -1, ErrMsg: "error :server name"})
		return
	}
	// fmt.Printf("pro_send>>>netInfo:%+v\n", netInfo)
	ind, ok := n.appNames.isAppNameIn(appname)
	if !ok {
		c.JSON(200, &NetRspInfo{ErrCode: -1, ErrMsg: "error :appname"})
		return
	}
	m := n.appNames.getMap(ind)
	val1, ok1 := m.Load(sid)
	// fmt.Printf("pro_send>>>ch1, ok1:%+v,%v\n", ch1, ok1)
	if ok1 {
		mInfo := val1.(*MemoryInfo)
		mInfo.ch <- &netInfo
	} else {
		c.JSON(200, &NetRspInfo{ErrCode: -1, ErrMsg: "error :sessid not exit"})
	}
	c.JSON(200, &NetRspInfo{ErrCode: 0, ErrMsg: "ok"})
}

type ServerConn struct {
	n       *NetServer
	ch      chan *NetReqInfo
	urlsend string
	client  *http.Client
	appname string
	m       *sync.Map
	sessid  string
	tk      *time.Timer
	timeout int64
}

func (n *NetServer) NewServerConn(appname, sessid string) (*ServerConn, error) {
	ind, ok := n.appNames.isAppNameIn(appname)
	if !ok {
		return nil, errors.New("appname not exist")
	}
	m := n.appNames.getMap(ind)
	t := time.NewTimer(300 * time.Second)
	defer func() { t.Stop() }()
	var v interface{}
loop:
	for {
		select {
		case <-t.C:
			return nil, errors.New("timeout")
		default:
			v, ok = m.Load(sessid)
			if ok {
				break loop
			}
			continue
		}
	}
	mInfo := v.(*MemoryInfo)
	client := &http.Client{Transport: http.DefaultTransport}
	urlSend := "http://" + mInfo.remoteIp + ":" + mInfo.remoteport + SENDPATH
	return &ServerConn{
		n:       n,
		ch:      mInfo.ch,
		urlsend: urlSend,
		client:  client,
		appname: appname,
		m:       m,
		sessid:  sessid,
		tk:      time.NewTimer(300 * time.Second),
		timeout: 300}, nil
}
func (s *ServerConn) Close() {
	s.tk.Stop()
	close(s.ch)
	s.m.Delete(s.sessid)
	// fmt.Printf("server conn close\n")
}
func (s *ServerConn) Recv() ([]byte, error) {
	// t := time.NewTimer(5 * time.Second)
	s.tk.Reset(time.Duration(s.timeout) * time.Second)
	for {
		select {
		case <-s.tk.C:
			return nil, errors.New("timeout")
		case data := <-s.ch:
			return data.Data, nil
		}
	}
	// t.Stop()
	// return nil, nil
}
func (s *ServerConn) Send(data []byte) error {
	reqInfo := NetReqInfo{data}
	res, _ := json.Marshal(&reqInfo)
	request, err := http.NewRequest("POST", s.urlsend, bytes.NewReader(res))
	if err != nil {
		return err
	}
	request.Header.Set(SESSID, s.sessid)
	request.Header.Set(UniqueServerName, s.n.uniqueServerName)
	rsp, err := s.client.Do(request)
	if err != nil {
		return err
	}
	defer func() { rsp.Body.Close() }()
	out, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	var rspInfo NetRspInfo
	if err := json.Unmarshal(out, &rspInfo); err != nil {
		return err
	}
	if rspInfo.ErrCode != 0 {
		return errors.New(rspInfo.ErrMsg)
	}
	return nil
}
