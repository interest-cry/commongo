package network

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type GinSender interface {
	WriteMessage(req requestForSend)
}
type requestForSend struct {
	NetworkType string `json:"network_type"`
	RemoteNid   string `json:"remote_nid"`
	Uid         string `json:"uid"`
	Key         string `json:"key"`
	Data        []byte `json:"data"`
}

type GinSenderMap struct {
	ginSenders sync.Map
}

var DefaultGinSenderMap = &GinSenderMap{
	ginSenders: sync.Map{},
	//ginSenders: map[string]GinSender{}
}

func NewGinSenderMap(timeout int) *GinSenderMap {
	ah := new(GinSenderMap)
	//ah.ginSenders = map[string]GinSender{}
	ah.ginSenders.Store(CHANCONN, NewEventBus(timeout))
	ah.ginSenders.Store(CACHECONN, NewHttpBigCache(timeout))
	//ah.ginSenders[CACHECONN] = NewHttpBigCache(timeout)
	//ah.ginSenders[CHANCONN] = NewEventBus(timeout)
	//fmt.Printf("=========>>>map:%+v\n", ah.ginSenders)
	return ah
}
func (m *GinSenderMap) GinSenderHandlerFunc(c *gin.Context) {
	var req requestForSend
	err := c.BindJSON(&req)
	if err != nil {
		DeLog.Infof(INFOPREFIX+"GinSenderFunc BindJson error:%v", err)
		return
	}
	h, ok := m.ginSenders.Load(req.NetworkType)
	//h, ok := a.ginSenders[req.NetworkType]
	if !ok {
		DeLog.Infof(INFOPREFIX + "GinSenderFunc get handler by network type")
		return
	}
	hand := h.(GinSender)
	hand.WriteMessage(req)
	return
}

func (m *GinSenderMap) GetGinSender(networkType string) GinSender {
	h, ok := m.ginSenders.Load(networkType)
	if !ok {
		panic("ginSenders.Load not ok")
	}
	ginSender, ok := h.(GinSender)
	if !ok {
		panic("h.(GinSender) not ok")
	}
	return ginSender
}

func (m *GinSenderMap) GetHttpBigCache() *HttpBigCache {
	h, ok := m.ginSenders.Load(CACHECONN)
	if !ok {
		panic("GetHttpBigCache error ,not ok")
	}
	httpBigCache, ok := h.(*HttpBigCache)
	if !ok {
		panic("GetHttpBigCache error ,not ok")
	}
	return httpBigCache
}
func (m *GinSenderMap) GetEventBus() *EventBus {
	h, ok := m.ginSenders.Load(CHANCONN)
	if !ok {
		panic("GetEventBus error,not ok")
	}
	eventBus, ok := h.(*EventBus)
	if !ok {
		panic("GetEventBus error,not ok")
	}
	return eventBus
}
