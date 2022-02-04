package network

//消息发送接口
type Messager interface {
	SendData(key string, val []byte) (int, error)
	RecvData(key string) ([]byte, error)
	Close() error
}
type Option func(o *Options)

const (
	TCPCONN  = "TCPCONN"
	HTTP     = "HTTP"
	HTTPCONN = "HTTPCONN"
	CHANCONN = "CHANCONN"
)

func NewMessager(netWorkType string, opts ...Option) (Messager, error) {
	//o := newOptions(opts...)
	switch netWorkType {
	case HTTP:
		break
	case HTTPCONN:
		return newHttpConn(opts...)

	case TCPCONN:
		return newTcpConn(opts...)

	case CHANCONN:
		return newChanConn(opts...)
	default:
		return newTcpConn(opts...)
	}
	return nil, nil
}
