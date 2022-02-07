package network

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type SendHandler interface {
	HandleMessage(req RequestForSend)
}
type RequestForSend struct {
	NetworkType string `json:"network_type"`
	RemoteNid   string `json:"remote_nid"`
	Uid         string `json:"uid"`
	Key         string `json:"key"`
	Data        []byte `json:"data"`
}

type AllHandler struct {
	handlers sync.Map
	//handlers map[string]SendHandler
}

var DefaultAllHandler = &AllHandler{
	handlers: sync.Map{},
	//handlers: map[string]SendHandler{}
}

func NewAllHandler(timeout int) *AllHandler {
	ah := new(AllHandler)
	//ah.handlers = map[string]SendHandler{}
	ah.handlers.Store(CHANCONN, NewEventBus(timeout))
	ah.handlers.Store(CACHECONN, NewHttpBigCache(timeout))
	//ah.handlers[CACHECONN] = NewHttpBigCache(timeout)
	//ah.handlers[CHANCONN] = NewEventBus(timeout)
	//fmt.Printf("=========>>>map:%+v\n", ah.handlers)
	return ah
}
func (a *AllHandler) SendHandlerFunc(c *gin.Context) {
	var req RequestForSend
	err := c.BindJSON(&req)
	if err != nil {
		DeLog.Infof(INFOPREFIX+"SendHandlerFunc BindJson error:%v", err)
		return
	}
	h, ok := a.handlers.Load(req.NetworkType)
	//h, ok := a.handlers[req.NetworkType]
	if !ok {
		DeLog.Infof(INFOPREFIX + "SendHandlerFunc get handler by network type")
		return
	}
	hand := h.(SendHandler)
	hand.HandleMessage(req)
	return
}

func (a *AllHandler) GetHandler(networkType string) SendHandler {
	h, ok := a.handlers.Load(networkType)
	//fmt.Printf("=========map:%+v\n", a.handlers)
	//h, ok := a.handlers[networkType]
	if !ok {
		DeLog.Infof(INFOPREFIX + "GetHandler get handler by network type")
		return nil
	}
	hand := h.(SendHandler)
	return hand
}
