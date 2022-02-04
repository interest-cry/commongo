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
func TestTcpConn_TimeOut(t *testing.T) {
	for i := 10; i < 100; i++ {
		o := newOptions(TimeOut(i))
		assert.Equal(t, i, o.TimeOut)
	}
}
