package eventbus

import (
	"errors"
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
	//rm          sync.RWMutex
	//subscribers sync.Map
	tick *time.Ticker
}

//订阅主题
func (eb *EventBus) Subscribe(topic string, ch chan DataEvent) {
	//eb.rm.Lock()
	//defer func() {
	//	eb.rm.Unlock()
	//}()
	//_, ok := eb.subscribers.Load(topic)
	_, ok := eb.subscribers[topic]
	if !ok {
		//eb.subscribers.Store(topic, ch)
		eb.subscribers[topic] = ch
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
			//ch, ok := eb.subscribers.Load(event.Topic)
			ch, ok := eb.subscribers[event.Topic]
			if ok {
				//这样做是因为切片引用相同的数组，即使他们是按值传递的
				//因此我们正在使用我们的元素创建一个新的切片，从而能正确的保持锁定
				//ch1, ok := ch.(chan DataEvent)
				//ch.(chan DataEvent) <- event
				//if !ok {
				//	return errors.New("类型断言出错")
				//}
				//ch1 <- event
				ch <- event
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
	//ch, ok := DefaultEventBus.subscribers.Load(topic)
	ch, ok := DefaultEventBus.subscribers[topic]
	if ok {
		//ch1, ok := ch.(chan DataEvent)
		//ch1, ok := ch.(chan DataEvent)
		if ok {
			//close(ch1)
			close(ch)
		}
	}
}

var DefaultEventBus = &EventBus{
	subscribers: map[string]DataEventChan{},
	//subscribers: sync.Map{},
	tick: time.NewTicker(1200 * time.Second),
}

func init() {
	//maps:=sync.Map{}
	//maps.
}
