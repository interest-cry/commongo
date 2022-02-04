package eventdriver

import "sync"

type DataEvent struct {
	Data  interface{}
	Topic string
}
type DataChannel chan DataEvent
type DataChannelSlice []DataChannel

//事件总线
type EventBus struct {
	subscribers map[string][]DataChannel
	rm          sync.RWMutex
}

//订阅主题
func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	defer func() {
		eb.rm.Unlock()
	}()
	prev, ok := eb.subscribers[topic]
	if ok {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
}

//发布主题
func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.Lock()
	defer func() {
		eb.rm.Unlock()
	}()
	chans, ok := eb.subscribers[topic]
	if ok {
		//这样做是因为切片引用相同的数组，即使他们是按值传递的
		//因此我们正在使用我们的元素创建一个新的切片，从而能正确的保持锁定
		chansTmp := append([]DataChannel{}, chans...)
		go func(data DataEvent, dataChanSlice []DataChannel) {
			for _, ch := range dataChanSlice {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, chansTmp)
	}
}

var eb = &EventBus{
	subscribers: map[string][]DataChannel{},
}
