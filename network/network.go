package network

import "errors"

//消息发送接口
type Messager interface {
	SendData(key string, val []byte) (int, error)
	RecvData(key string) ([]byte, error)
	Close() error
}
type Option func(o *Options)

const (
	TCPCONN   = "TCPCONN"
	HTTP      = "HTTP"
	CACHECONN = "CACHECONN"
	CHANCONN  = "CHANCONN"
)

type newConnFunc func(opts ...Option) (Messager, error)

var NetworkRegister []string = []string{TCPCONN, CACHECONN, CHANCONN}
var NetworkConnRegister map[string]newConnFunc = map[string]newConnFunc{
	CHANCONN:  newChanConn,
	CACHECONN: newCacheConn,
	TCPCONN:   newTcpConn,
}
var NetworkMap map[string]string = map[string]string{
	"tcp":   TCPCONN,
	"http":  HTTP,
	"cache": CACHECONN,
	"chan":  CHANCONN,
}

func NewMessager(netWorkType string, opts ...Option) (Messager, error) {
	//o := newOptions(opts...)
	opts = append(opts, NetworkType(netWorkType))
	DeLog.Infof(INFOPREFIX+"netWorkType:%v", netWorkType)
	connFunc, ok := NetworkConnRegister[netWorkType]
	if !ok {
		return nil, errors.New("error:not in NetworkConnRegister")
	}
	return connFunc(opts...)
}
