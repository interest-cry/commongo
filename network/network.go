package network

import "errors"

//消息发送、接收接口
type Communicator interface {
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

type newConnFunc func(opts ...Option) (Communicator, error)

var NetworkRegister []string = []string{TCPCONN, CACHECONN, CHANCONN}
var CommunicatorRegister map[string]newConnFunc = map[string]newConnFunc{
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

func NewCommunicator(netWorkType string, opts ...Option) (Communicator, error) {
	//o := newOptions(opts...)
	opts = append(opts, NetworkType(netWorkType))
	DeLog.Infof(INFOPREFIX+"netWorkType:%v", netWorkType)
	connFunc, ok := CommunicatorRegister[netWorkType]
	if !ok {
		return nil, errors.New("error:not in NetworkConnRegister")
	}
	return connFunc(opts...)
}
