package httpnetgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type NetClient struct {
	ip        string
	port      string
	engine    *gin.Engine
	cliMemory *sync.Map
}

func NewNetClient(ipAnPort string) *NetClient {
	ret := strings.Split(ipAnPort, ":")
	if len(ret) != 2 {
		panic("error:ip and port is wrong.please use ip:port(127.0.0.1:18888)")
	}
	var n = NetClient{
		ip:        ret[0],
		port:      ret[1],
		engine:    gin.New(),
		cliMemory: &sync.Map{}}
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

func (n *NetClient) sendcall(c *gin.Context) {
	//n.engine.POST("/send", func(c *gin.Context) {
	var netInfo NetReqInfo
	c.ShouldBindJSON(&netInfo)
	sid := c.GetHeader(SESSID)
	sername := c.GetHeader(UniqueServerName)
	uid := sername + "_" + sid
	// fmt.Printf("=======client uid:%+v\n", uid)
	val1, ok1 := n.cliMemory.Load(uid)
	// fmt.Printf("pro_send>>>ch1, ok1:%+v,%v\n", ch1, ok1)
	if ok1 {
		ch := val1.(chan *NetReqInfo)
		ch <- &netInfo
	} else {
		c.JSON(200, &NetRspInfo{ErrCode: -1, ErrMsg: "conn not create"})
		return
	}
	c.JSON(200, &NetRspInfo{ErrCode: 0, ErrMsg: "ok"})
	//})
}

type ClientConn struct {
	n       *NetClient
	ch      chan *NetReqInfo
	urlsend string
	client  *http.Client
	appname string
	sessid  string
	sername string
	uid     string
	tk      *time.Timer
	timeout int64
}

func (n *NetClient) NewClientConn(remoteIpAndPort string, appname string) (*ClientConn, error) {
	client := &http.Client{Transport: http.DefaultTransport}
	url := "http://" + remoteIpAndPort + INITCONNPATH
	req, err := http.NewRequest("GET", url, bytes.NewReader(nil))
	if err != nil {
		return nil, err
	}
	req.Header.Set(APPNAME, appname)
	req.Header.Set(REMOTEIP, n.ip)
	req.Header.Set(REMOTEPORT, n.port)
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { rsp.Body.Close() }()
	out, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	var netRsp NetRspInfo
	if err := json.Unmarshal(out, &netRsp); err != nil {
		return nil, err
	}
	if netRsp.ErrCode != 0 {
		return nil, errors.New("init conn failed")
	}
	sername := rsp.Header.Get(UniqueServerName)
	sid := rsp.Header.Get(SESSID)
	uid := sername + "_" + sid
	ch := make(chan *NetReqInfo)
	n.cliMemory.Store(uid, ch)
	urlsend := "http://" + remoteIpAndPort + SENDPATH
	con := ClientConn{
		n:       n,
		ch:      ch,
		urlsend: urlsend,
		client:  client,
		appname: appname,
		sessid:  sid,
		sername: sername,
		uid:     uid,
		tk:      time.NewTimer(300 * time.Second),
		timeout: 300}
	return &con, nil
}
func (c *ClientConn) Close() {
	// v, ok := c.n.cliMemory.Load(c.uid)
	close(c.ch)
	c.n.cliMemory.Delete(c.uid)
	c.tk.Stop()
	// fmt.Printf("client conn close\n")
}
func (c *ClientConn) Recv() ([]byte, error) {
	// t := time.NewTimer(5 * time.Second)
	c.tk.Reset(time.Duration(c.timeout) * time.Second)
	for {
		select {
		case <-c.tk.C:
			return nil, errors.New("timeout")
		case data := <-c.ch:
			return data.Data, nil
		}
	}
	// t.Stop()
	// return nil, nil
}
func (c *ClientConn) Send(data []byte) error {
	reqInfo := NetReqInfo{data}
	res, _ := json.Marshal(&reqInfo)
	request, err := http.NewRequest("POST", c.urlsend, bytes.NewReader(res))
	if err != nil {
		return err
	}
	request.Header.Set(APPNAME, c.appname)
	request.Header.Set(SESSID, c.sessid)
	request.Header.Set(UniqueServerName, c.sername)
	rsp, err := c.client.Do(request)
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

func (c *ClientConn) SendStart(ipAndPort, appname string) error {
	cli := http.Client{Transport: http.DefaultTransport}
	// serverIp := "127.0.0.1"
	rq, err := http.NewRequest("GET", "http://"+ipAndPort+"/"+appname, nil)
	if err != nil {
		return err
	}
	rq.Header.Set(APPNAME, appname)
	rq.Header.Set(SESSID, c.sessid)
	rs, err := cli.Do(rq)
	if err != nil {
		return err
	}
	defer func() {
		rs.Body.Close()
	}()
	return nil
}
func (c *ClientConn) GetSessid() string {
	return c.sessid
}
