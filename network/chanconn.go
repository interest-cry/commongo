package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type DataEvent struct {
	Key   string
	Data  []byte
	Topic string
}
type DataEventChan chan DataEvent

//事件总线
type EventBus struct {
	//subscribers map[string]DataEventChan
	//rm          sync.RWMutex
	subscribers sync.Map
	timeout     time.Duration
	//tick        *time.Timer
}

//订阅主题
func (eb *EventBus) Subscribe(topic string, ch chan DataEvent) {
	//eb.rm.Lock()
	//defer func() {
	//	eb.rm.Unlock()
	//}()
	_, ok := eb.subscribers.Load(topic)
	//_, ok := eb.subscribers[topic]
	if !ok {
		eb.subscribers.Store(topic, ch)
		//eb.subscribers[topic] = ch
	}
}

//发布主题
func (eb *EventBus) Publish(event DataEvent) error {
	//eb.rm.Lock()
	//defer func() {
	//	eb.rm.Unlock()
	//}()
	tick := time.NewTimer(eb.timeout)
	defer func() {
		tick.Stop()
	}()
	for {
		select {
		case <-tick.C:
			return errors.New("event bus publish timeout")
		default:
			ch, ok := eb.subscribers.Load(event.Topic)
			//ch, ok := eb.subscribers[event.Topic]
			if ok {
				ch1, ok := ch.(chan DataEvent)
				//ch.(chan DataEvent) <- event
				if ok {
					ch1 <- event
				} else {
					return errors.New("类型断言出错")
				}
				//ch <- event
				return nil
			} else {
				continue
			}
		}
	}
	return nil
}
func (eb *EventBus) Close(topic string) {
	ch, ok := eb.subscribers.Load(topic)
	DeLog.Infof(INFOPREFIX+"Close,ch:%v,ok:%v", ch, ok)
	//ch, ok := DefaultEventBus.subscribers[topic]
	if ok {
		ch1, ok := ch.(chan DataEvent)
		if ok {
			close(ch1)
			//DeLog.Infof(INFOPREFIX + "close chan")
			//close(ch)
		}
	}
}

var DefaultEventBus = &EventBus{
	//subscribers: map[string]DataEventChan{},
	subscribers: sync.Map{},
	timeout:     1200 * time.Second,
}

func NewEventBus(timeout int) *EventBus {
	return &EventBus{
		subscribers: sync.Map{},
		timeout:     time.Duration(timeout) * time.Second,
	}
}

func (e *EventBus) WriteMessageGin(c *gin.Context) {
	var req requestForSend
	err := c.BindJSON(&req)
	if err != nil {
		DeLog.Infof(INFOPREFIX+"SaveData BindJson error:%v", err)
		return
	}
	//err = hb.bigC.Set(req.Key, req.Data)
	topic := req.RemoteNid + "_" + req.Uid
	err = e.Publish(DataEvent{Topic: topic, Key: req.Key, Data: req.Data})
	if err != nil {
		DeLog.Infof(INFOPREFIX+"SaveData set val error:%v", err)
		return
	}
	//DeLog.Infof(INFOPREFIX + "save data ok")
	return
	//message.Log.Infof("===>>本地缓存 set data ok")
	//todo:不需要响应返回
	//c.JSON(200, gin.H{
	//	"msg": "ok",
	//})
}

func (e *EventBus) WriteMessage(req requestForSend) {
	topic := req.RemoteNid + "_" + req.Uid
	err := e.Publish(DataEvent{Topic: topic, Key: req.Key, Data: req.Data})
	if err != nil {
		DeLog.Infof(INFOPREFIX+"SaveData set val error:%v", err)
		return
	}
	//DeLog.Infof(INFOPREFIX + "save data ok")
	return
	//message.Log.Infof("===>>本地缓存 set data ok")
	//todo:不需要响应返回
	//c.JSON(200, gin.H{
	//	"msg": "ok",
	//})
}

type ChanConn struct {
	o          *Options
	eventBus   *EventBus
	httpClient *http.Client
	ch         chan DataEvent
	tick       *time.Timer
	timeout    time.Duration
}

func newChanConn(opts ...Option) (Communicator, error) {
	o := newOptions(opts...)
	topic := o.RemoteNid + "_" + o.Uid
	//关键,订阅
	ch := make(chan DataEvent)
	o.EventB.Subscribe(topic, ch)
	timeout := time.Duration(o.TimeOut) * time.Second
	return &ChanConn{
		o:          o,
		eventBus:   o.EventB,
		httpClient: &http.Client{Transport: http.DefaultTransport},
		ch:         ch,
		tick:       time.NewTimer(timeout),
		timeout:    timeout,
	}, nil
}

func (c *ChanConn) SendData(key string, val []byte) (int, error) {
	req := requestForSend{
		NetworkType: c.o.NetworkType,
		RemoteNid:   c.o.LocalNid,
		Uid:         c.o.Uid,
		Key:         key,
		Data:        val,
	}
	dataJson, _ := json.Marshal(&req)
	rsp, err := c.httpClient.Post(c.o.SendUrl, "application/json", bytes.NewReader(dataJson))
	if err != nil {
		return 0, err
	}
	rsp.Body.Close()
	return len(val), nil
}
func (c *ChanConn) RecvData(key string) ([]byte, error) {
	//todo:
	dataEvent, ok := <-c.ch
	if ok {
		if dataEvent.Key == key {
			return dataEvent.Data, nil
		}
		return nil, errors.New("key error")
	}
	return nil, errors.New("chan closed")
}
func (c *ChanConn) Close() error {
	c.tick.Stop()
	return nil
}
