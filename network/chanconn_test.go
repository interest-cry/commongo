package network

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	guuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestEventBus_NewEventBus(t *testing.T) {
	eventBus := NewEventBus(1200)
	DeLog.Infof(INFOPREFIX+"eventBus:%+v", eventBus)
	DeLog.Infof(INFOPREFIX+"DefaultEventBus:%+v", DefaultEventBus)
}
func TestEventBus_Publish(t *testing.T) {
	eventBus := NewEventBus(1)
	topic := "topic1"
	key := "key1"
	data := []byte("333")
	event := DataEvent{
		Key:   key,
		Data:  data,
		Topic: topic,
	}
	err := eventBus.Publish(event)
	DeLog.Infof(INFOPREFIX+"timeout,error:%v", err)
}
func TestEventBus_Publish1(t *testing.T) {
	eventBus := NewEventBus(5)
	topic := "topic1"
	key := "key1"
	data := []byte("333")
	event := DataEvent{
		Key:   key,
		Data:  data,
		Topic: topic,
	}
	//先订阅主题
	ch := make(chan DataEvent)
	eventBus.Subscribe(topic, ch)
	go func() {
		err := eventBus.Publish(event)
		//DeLog.Infof(INFOPREFIX+"error:%v", err)
		assert.NoError(t, err)
		eventBus.Close(topic)
	}()
	//time.Sleep(time.Second * 5)
	//消费
	ret, ok := <-ch
	assert.True(t, true, ok)
	DeLog.Infof(INFOPREFIX+"ret:%v", ret)
	ret, ok = <-ch
	assert.False(t, false, ok)
	DeLog.Infof(INFOPREFIX+"ret:%v,ok:%+v", ret, ok)
}

func TestEventBus_Publish2(t *testing.T) {
	eventBus := NewEventBus(5)
	topic := "topic1"
	key := "key1"
	data := []byte("333")
	event := DataEvent{
		Key:   key,
		Data:  data,
		Topic: topic,
	}
	//订阅主题
	ch := make(chan DataEvent)
	//eventBus.Subscribe(topic, ch)
	//先发布
	go func() {
		err := eventBus.Publish(event)
		//DeLog.Infof(INFOPREFIX+"error:%v", err)
		assert.NoError(t, err)
		eventBus.Close(topic)
	}()
	time.Sleep(time.Second * 2)
	eventBus.Subscribe(topic, ch)
	//消费
	ret, ok := <-ch
	assert.True(t, true, ok)
	DeLog.Infof(INFOPREFIX+"ret:%v", ret)
	ret, ok = <-ch
	assert.False(t, false, ok)
	DeLog.Infof(INFOPREFIX+"ret:%v,ok:%+v", ret, ok)
}
func TestEventBus_Close(t *testing.T) {

}
func TestEventBus_Subscribe(t *testing.T) {

}

type mockChanServer struct {
	eventBus *EventBus
	r        *gin.Engine
}

func newMockChanServer(timeout int) *mockChanServer {
	return &mockChanServer{
		eventBus: NewEventBus(timeout),
		r:        gin.New(),
	}
}
func (s *mockChanServer) addPath(relativePath string) {
	s.r.POST(relativePath, s.eventBus.EventBusHandlerFunc)
}

func TestEventBus_EventBusHandlerFunc(t *testing.T) {
	relativePath := "/v1/send"
	chanServer := newMockChanServer(5)
	chanServer.addPath(relativePath)
	//启动一个服务
	service := httptest.NewServer(chanServer.r)
	defer func() { service.Close() }()
	uid := "20220205"
	nid0 := "nid0"
	nid1 := "nid1"
	datasetNum := 10
	wg := sync.WaitGroup{}
	wg.Add(2)
	//接收数据,服务启动的节点nid0
	go func() {
		defer func() { wg.Done() }()
		//创建一个连接，订阅消息
		chanconn, err := newChanConn(EventBusSet(chanServer.eventBus),
			Uid(uid), LocalNid(nid0), RemoteNid(nid1))
		assert.NoError(t, err)
		defer func() {
			chanconn.Close()
		}()
		for i := 0; i < datasetNum; i++ {
			key := "key_" + strconv.Itoa(i)
			data, err := chanconn.RecvData(key)
			assert.NoError(t, err)
			DeLog.Infof("RecvData,key[%v],data[%v]", key, string(data))
		}
	}()
	//客户端发送数据nid1
	go func() {
		defer func() { wg.Done() }()
		//创建一个连接，订阅消息
		chanconn, err := newChanConn(Uid(uid), LocalNid(nid1),
			RemoteNid(nid0), SendUrl(service.URL+relativePath))
		assert.NoError(t, err)
		defer func() {
			chanconn.Close()
		}()
		for i := 0; i < datasetNum; i++ {
			key := "key_" + strconv.Itoa(i)
			data := []byte("msg_" + key)
			_, err := chanconn.SendData(key, data)
			assert.NoError(t, err)
			//DeLog.Infof("SendData,key[%v],data[%v]", key, string(data))
		}
	}()
	wg.Wait()
}

func TestEventBus_EventBusHandlerFuncBothWay(t *testing.T) {
	//启动第一个服务
	Nid1 := "nid1"
	relativePath1 := "/v1/send1"
	chanServer1 := newMockChanServer(5)
	chanServer1.addPath(relativePath1)
	service1 := httptest.NewServer(chanServer1.r)
	defer func() { service1.Close() }()
	//启动第一个服务,END
	//启动第二个服务
	Nid2 := "nid2"
	relativePath2 := "/v1/send2"
	chanServer2 := newMockChanServer(5)
	chanServer2.addPath(relativePath2)
	service2 := httptest.NewServer(chanServer2.r)
	defer func() { service2.Close() }()
	//启动第二个服务,END
	uid := "20220205"
	datasetNum := 10
	wg := sync.WaitGroup{}
	wg.Add(2)
	//第一个服务节点
	go func() {
		defer func() { wg.Done() }()
		//创建一个通信连接，订阅消息
		chanconn, err := newChanConn(EventBusSet(chanServer1.eventBus),
			Uid(uid), LocalNid(Nid1), RemoteNid(Nid2),
			SendUrl(service2.URL+relativePath2))
		assert.NoError(t, err)
		defer func() {
			chanconn.Close()
		}()
		for i := 0; i < datasetNum; i++ {
			//先接收
			key := Nid2 + "_" + strconv.Itoa(i)
			data, err := chanconn.RecvData(key)
			assert.NoError(t, err)
			DeLog.Infof("[%v]RecvData,key[%v],data[%v]", Nid1, key, string(data))
			//再发送
			send_key := Nid1 + "_" + strconv.Itoa(i)
			_, err = chanconn.SendData(send_key, []byte("msg_"+send_key))
			assert.NoError(t, err)
		}
	}()
	////第二个服务节点
	go func() {
		defer func() { wg.Done() }()
		//创建一个连接，订阅消息
		chanconn, err := newChanConn(EventBusSet(chanServer2.eventBus),
			Uid(uid), LocalNid(Nid2), RemoteNid(Nid1),
			SendUrl(service1.URL+relativePath1))
		assert.NoError(t, err)
		defer func() {
			chanconn.Close()
		}()
		for i := 0; i < datasetNum; i++ {
			//先发送
			send_key := Nid2 + "_" + strconv.Itoa(i)
			//data := []byte("msg_" + send_key)
			_, err := chanconn.SendData(send_key, []byte("msg_"+send_key))
			assert.NoError(t, err)
			//DeLog.Infof("SendData,key[%v],data[%v]", key, string(data))
			//再接收
			key := Nid1 + "_" + strconv.Itoa(i)
			data, err := chanconn.RecvData(key)
			assert.NoError(t, err)
			DeLog.Infof("[%v]RecvData,key[%v],data[%v]", Nid2, key, string(data))
		}
	}()
	wg.Wait()
}

func TestEventBus_EventBusHandlerFuncThreeNode(t *testing.T) {
	//启动第一个服务
	Nid1 := "nid1"
	relativePath1 := "/v1/send1"
	chanServer1 := newMockChanServer(5)
	chanServer1.addPath(relativePath1)
	service1 := httptest.NewServer(chanServer1.r)
	defer func() { service1.Close() }()
	//启动第一个服务,END
	//启动第二个服务
	Nid2 := "nid2"
	relativePath2 := "/v1/send2"
	chanServer2 := newMockChanServer(5)
	chanServer2.addPath(relativePath2)
	service2 := httptest.NewServer(chanServer2.r)
	defer func() { service2.Close() }()
	//启动第二个服务,END
	//启动第三个服务
	Nid3 := "nid3"
	relativePath3 := "/v1/send3"
	chanServer3 := newMockChanServer(5)
	chanServer3.addPath(relativePath3)
	service3 := httptest.NewServer(chanServer3.r)
	defer func() { service3.Close() }()
	//启动第三个服务,END
	uid := "20220205"
	uuid1, err := uuid.NewV4()
	assert.NoError(t, err)

	uid = uuid1.String()
	uid = guuid.New().String()
	datasetNum := 10000
	wg := sync.WaitGroup{}
	wg.Add(3)
	//第一个服务节点
	go func() {
		defer func() { wg.Done() }()
		//创建一个通信连接，订阅消息
		chanconn, err := newChanConn(EventBusSet(chanServer1.eventBus),
			Uid(uid), LocalNid(Nid1), RemoteNid(Nid2),
			SendUrl(service2.URL+relativePath2))
		assert.NoError(t, err)
		defer func() {
			chanconn.Close()
		}()
		for i := 0; i < datasetNum; i++ {
			//先接收
			key := Nid2 + "_" + strconv.Itoa(i)
			data, err := chanconn.RecvData(key)
			assert.NoError(t, err)
			DeLog.Infof("[%v]RecvData,key[%v],data[%v]", Nid1, key, string(data))
			//再发送
			send_key := Nid1 + "_" + strconv.Itoa(i)
			_, err = chanconn.SendData(send_key, []byte("msg_"+send_key))
			assert.NoError(t, err)
		}
		//1-->3
		chanconn3, err := newChanConn(EventBusSet(chanServer1.eventBus),
			Uid(uid), LocalNid(Nid1), RemoteNid(Nid3),
			SendUrl(service3.URL+relativePath3))
		assert.NoError(t, err)
		defer func() {
			chanconn3.Close()
		}()
		for i := 0; i < datasetNum; i++ {
			//发送
			send_key := Nid1 + "_" + strconv.Itoa(i)
			_, err = chanconn3.SendData(send_key, []byte("msg_"+send_key))
			assert.NoError(t, err)
		}
	}()
	//第二个服务节点
	go func() {
		defer func() { wg.Done() }()
		//创建一个连接，订阅消息
		chanconn, err := newChanConn(EventBusSet(chanServer2.eventBus),
			Uid(uid), LocalNid(Nid2), RemoteNid(Nid1),
			SendUrl(service1.URL+relativePath1))
		assert.NoError(t, err)
		defer func() {
			chanconn.Close()
		}()
		for i := 0; i < datasetNum; i++ {
			//先发送
			send_key := Nid2 + "_" + strconv.Itoa(i)
			//data := []byte("msg_" + send_key)
			_, err := chanconn.SendData(send_key, []byte("msg_"+send_key))
			assert.NoError(t, err)
			//DeLog.Infof("SendData,key[%v],data[%v]", key, string(data))
			//再接收
			key := Nid1 + "_" + strconv.Itoa(i)
			data, err := chanconn.RecvData(key)
			assert.NoError(t, err)
			DeLog.Infof("[%v]RecvData,key[%v],data[%v]", Nid2, key, string(data))
		}
		//2-->3
		chanconn3, err := newChanConn(EventBusSet(chanServer1.eventBus),
			Uid(uid), LocalNid(Nid2), RemoteNid(Nid3),
			SendUrl(service3.URL+relativePath3))
		assert.NoError(t, err)
		defer func() {
			chanconn3.Close()
		}()
		for i := 0; i < datasetNum; i++ {
			//发送
			send_key := Nid2 + "_" + strconv.Itoa(i)
			_, err = chanconn3.SendData(send_key, []byte("msg_"+send_key))
			assert.NoError(t, err)
		}
	}()
	//第三个服务节点
	go func() {
		defer func() { wg.Done() }()
		//1-->3
		chanconn1, err := newChanConn(EventBusSet(chanServer3.eventBus),
			Uid(uid), LocalNid(Nid3), RemoteNid(Nid1),
			SendUrl(service1.URL+relativePath1))
		assert.NoError(t, err)
		defer func() {
			chanconn1.Close()
		}()
		//2-->3
		chanconn2, err := newChanConn(EventBusSet(chanServer3.eventBus),
			Uid(uid), LocalNid(Nid3), RemoteNid(Nid2),
			SendUrl(service1.URL+relativePath1))
		assert.NoError(t, err)
		defer func() {
			chanconn2.Close()
		}()
		for i := 0; i < datasetNum; i++ {
			//接收
			key1 := Nid1 + "_" + strconv.Itoa(i)
			data1, err := chanconn1.RecvData(key1)
			assert.NoError(t, err)
			DeLog.Infof("[%v]RecvData,key[%v],data[%v]", Nid3, key1, string(data1))
			//
			key2 := Nid2 + "_" + strconv.Itoa(i)
			data2, err := chanconn2.RecvData(key2)
			assert.NoError(t, err)
			DeLog.Infof("[%v]RecvData,key[%v],data[%v]", Nid3, key2, string(data2))
		}
	}()
	wg.Wait()
}
