package network

//消息发送接口
type Messager interface {
	SendData(key string, val []byte) (int, error)
	RecvData(key string) ([]byte, error)
	Close() error
}
type Option func(o *Options)

const (
	TCP       = "TCP"
	HTTP      = "HTTP"
	HTTPCACHE = "HTTPCACHE"
)

func NewMessager(opts ...Option) (Messager, error) {
	o := newOptions(opts...)
	switch o.NetWorkType {
	case HTTP:
		break
	case HTTPCACHE:
		break
	default:
		//TCP
		return newTcpConn(o)
	}
	return nil, nil
}
