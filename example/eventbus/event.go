package eventbus

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
	subscribers map[string]DataEventChan
	rm          sync.RWMutex
	tick        *time.Ticker
}

//订阅主题
func (eb *EventBus) Subscribe(topic string, ch DataEventChan) {
	eb.rm.Lock()
	defer func() {
		eb.rm.Unlock()
	}()
	_, ok := eb.subscribers[topic]
	if !ok {
		eb.subscribers[topic] = ch
	}
}

//发布主题
func (eb *EventBus) Publish(event DataEvent) error {
	eb.rm.Lock()
	defer func() {
		eb.rm.Unlock()
	}()
	for {
		select {
		case <-eb.tick.C:
			return errors.New("event bus publish timeout")
		default:
			ch, ok := eb.subscribers[event.Topic]
			if ok {
				//这样做是因为切片引用相同的数组，即使他们是按值传递的
				//因此我们正在使用我们的元素创建一个新的切片，从而能正确的保持锁定
				ch <- event
				return nil
			} else {
				continue
			}
		}
	}
	return nil
}

var DefaultEventBus = &EventBus{
	subscribers: map[string]DataEventChan{},
	tick:        time.NewTicker(1200 * time.Second),
}
