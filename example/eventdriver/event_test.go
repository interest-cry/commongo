package eventdriver

import (
	"fmt"
	"testing"
)

func publishTo(topic string, data string) {
	for {
		eb.Publish(topic, data)
		//time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}
}
func printDataEvent(ch string, data DataEvent) {
	fmt.Printf("Channel:%s;Topic:%s;Data:%+v\n", ch, data.Topic, data.Data)
}
func TestEventBus_Publish(t *testing.T) {
	ch1 := make(chan DataEvent)
	ch2 := make(chan DataEvent)
	ch3 := make(chan DataEvent)
	eb.Subscribe("topic1", ch1)
	eb.Subscribe("topic2", ch2)
	eb.Subscribe("topic2", ch3)
	go publishTo("topic1", "hi topic 1")
	go publishTo("topic2", "welcome to topic 2")
	for {
		select {
		case d := <-ch1:
			go printDataEvent("ch1", d)
		case d := <-ch2:
			go printDataEvent("ch2", d)
		case d := <-ch3:
			go printDataEvent("ch3", d)
		}
		//time.Sleep(time.Second)
	}
}
