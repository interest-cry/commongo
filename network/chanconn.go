package network

import (
	"errors"
	"sync"
	"time"
)

type DataEvent struct {
	Key   string
	Data  []byte
	Topic string
}
type DataEventChan chan DataEvent

//type DataChannelSlice []DataChannel

//事件总线
type EventBus struct {
	//subscribers map[string]DataEventChan
	//rm          sync.RWMutex
	subscribers sync.Map
	tick        *time.Ticker
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
	for {
		select {
		case <-eb.tick.C:
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
	//eb.rm.Lock()
	//defer func() {
	//	eb.rm.Unlock()
	//}()
	ch, ok := DefaultEventBus.subscribers.Load(topic)
	//ch, ok := DefaultEventBus.subscribers[topic]
	if ok {
		ch1, ok := ch.(chan DataEvent)
		if ok {
			close(ch1)
			//close(ch)
		}
	}
}

var DefaultEventBus = &EventBus{
	//subscribers: map[string]DataEventChan{},
	subscribers: sync.Map{},
	tick:        time.NewTicker(1200 * time.Second),
}

func init() {
	//maps:=sync.Map{}
	//maps.
}

type ChanConn struct {
	eventBus *EventBus
}

func EventBusSet(eventB *EventBus) Option {
	return func(o *Options) {
		o.EventB = eventB
	}
}
func newChanConn(opts ...Option) (*ChanConn, error) {
	o := newOptions(opts...)
	return &ChanConn{
		eventBus: o.EventB}, nil
}

func (c *ChanConn) SendData(key string, val []byte) (int, error) {
	return len(val), nil
}
func (c *ChanConn) RecvData(key string) ([]byte, error) {
	return nil, nil
}
func (c *ChanConn) Close() error {
	return nil
}
