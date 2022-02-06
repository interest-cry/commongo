package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/***
go test -v -run TestIp
go test -v -test.run=TestIp
go test -v ./*.go
go test -v ./*.go -test.run=.
go test -v ./*.go -test.bench=.
go test -v options_test.go -test.run TestIp
***/
/**
会出现未定义的情况, 这是因为定义在其他文件里, 需要加上定义的文件.
go test -v options.go options_test.go
go test -v ./*.go -test.run=. 全测
go test -v ./*.go -test.run=TestTcpConn_NewTcpConn
**/
//go test -v hello.go hello_test.go
//go test -v  -test.run="TestA*";加前缀
/***
测试覆盖率
go test -v -cover
go test -v -coverprofile=a.out -test.run="TestA*" # 把测试结果保存在 a.out
go tool cover -html=./a.out  # 通过浏览器打开, 可以看到覆盖经过的函数
 ****/
func Test_Ip(t *testing.T) {
	ip := "192.168.1.6"
	o := newOptions(Ip(ip))
	assert.True(t, o != nil, "o==nil")
	assert.Equal(t, ip, o.Ip)
}

func Test_Port(t *testing.T) {
	port := 3000
	o := newOptions(Port(port))
	assert.True(t, o != nil, "o==nil")
	assert.Equal(t, port, o.Port)
}

//func TestNetWorkType(t *testing.T) {
//	netTypes := []string{TCP, HTTP, HTTPCACHE}
//	for i, v := range netTypes {
//		assert.Equal(t, netTypes[i], v)
//		o := newOptions(NetWorkType(v))
//		assert.Equal(t, v, o.NetWorkType)
//	}
//}

func Test_TimeOut(t *testing.T) {
	for i := 10; i < 100; i++ {
		o := newOptions(TimeOut(i))
		assert.Equal(t, i, o.TimeOut)
	}
}

func Test_ClientOrServer(t *testing.T) {
	cs := []string{CLIENT, SERVER}
	o := newOptions()
	assert.Equal(t, o.ClientOrServer, CLIENT)
	for i, v := range cs {
		assert.Equal(t, cs[i], v)
		o := newOptions(ClientOrServer(v))
		assert.Equal(t, v, o.ClientOrServer)
	}
}
func Benchmark_Ip(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Ip("127.0.0.1")
	}
}
