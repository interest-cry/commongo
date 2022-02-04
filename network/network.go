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

func NewMessager(netWorkType string, opts ...Option) (Messager, error) {
	//o := newOptions(opts...)
	switch netWorkType {
	case HTTP:
		break
	case HTTPCACHE:
		return newHttpConn(opts...)
		break
	default:
		//TCP
		return newTcpConn(opts...)
	}
	return nil, nil
}
