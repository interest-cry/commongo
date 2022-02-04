package eventbus

import (
	"fmt"
	"github.com/Zhenhanyijiu/commongo/example/server"
	"strconv"
	"testing"
)

func publishTo(topic string) {
	i := 0
	datasetNum := 977 * 5000
	datasetNum = 1000000
	dataSrcLen := 102400
	srcData, _ := server.GenRandDataDebug(11, datasetNum, dataSrcLen)
	//defer func() {
	//	DefaultEventBus.Close(topic)
	//}()
	for i < datasetNum {
		key := "key_" + strconv.Itoa(i)
		data := append([]byte{}, srcData...)
		data = []byte(key)
		event := DataEvent{Topic: topic, Data: data, Key: key}
		err := DefaultEventBus.Publish(event)
		if err != nil {
			panic(err)
		}
		//time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		i++
	}
}
func printDataEvent(ch string, event DataEvent) {
	fmt.Printf("Channel:%s;Topic:%s;key:%+v,Data:%+v\n", ch, event.Topic, event.Key, string(event.Data))
}
func TestEventBus_Publish(t *testing.T) {
	ch1 := make(chan DataEvent)
	ch2 := make(chan DataEvent)
	//ch3 := make(chan DataEvent)
	//订阅主题
	DefaultEventBus.Subscribe("topic1", ch1)
	DefaultEventBus.Subscribe("topic2", ch2)
	//eb.Subscribe("topic2", ch2)
	//eb.Subscribe("topic2", ch3)
	go publishTo("topic1")
	go publishTo("topic2")
	//go publishTo("topic2", "welcome to topic 2")
	for {
		//select {
		//case d1, ok := <-ch1:
		//	if !ok {
		//		return
		//	}
		//	printDataEvent("ch1", d1)
		//case d2, ok := <-ch2:
		//	if !ok {
		//		return
		//	}
		//	printDataEvent("ch2", d2)
		//case d := <-ch3:
		//	go printDataEvent("ch3", d)
		//}
		//time.Sleep(time.Second)
		d1 := <-ch1
		printDataEvent("ch1", d1)
		d2 := <-ch2
		printDataEvent("ch2", d2)
	}
}
