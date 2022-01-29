package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

const (
	CLIENT = "CLIENT"
	SERVER = "SERVER"
)

type connResult struct {
	c   net.Conn
	err error
}
type TcpConn struct {
	o           *Options
	c           *net.TCPConn
	listener    net.Listener
	headSendBuf []byte
	headRecvBuf []byte
}

func newTcpConn(o *Options) (*TcpConn, error) {
	//o := newOptions(opts...)
	//dur := time.Tick(time.Second * time.Duration(o.TimeOut))
	//可以手动停止任务
	tick := time.NewTicker(time.Second * time.Duration(o.TimeOut))
	ip := o.Ip
	portInt := o.Port
	port := strconv.Itoa(portInt)
	//tcp server
	if o.ClientOrServer == SERVER {
		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			return nil, err
		}
		out := make(chan connResult, 1)
		go func() {
			defer func() {
				fmt.Printf("listener accept go routine exit ok\n")
			}()
			c, err := listener.Accept()
			out <- connResult{c, err}
		}()
		var connRet connResult
		var errRet error
		select {
		case connRet = <-out:
			errRet = connRet.err
		case <-tick.C:
			//超时
			//listener.Close()
			errRet = errors.New("server listen timeout")
		}
		//停止tick
		tick.Stop()
		if errRet != nil {
			listener.Close()
			return nil, errRet
		}
		fmt.Printf("server: accept ok,remoteaddr:%v\n", connRet.c.RemoteAddr())
		return &TcpConn{
			o:           o,
			c:           connRet.c.(*net.TCPConn),
			listener:    listener,
			headSendBuf: make([]byte, 4),
			headRecvBuf: make([]byte, 4)}, nil
	}
	//tcp client
	if o.ClientOrServer == CLIENT {
		for {
			select {
			case <-tick.C:
				tick.Stop()
				return nil, errors.New("dial timeout")
			default:
				c, err := net.Dial("tcp", ip+":"+port)
				if err != nil {
					continue
				}
				//释放资源
				tick.Stop()
				fmt.Printf("client: dial ok,remoteaddr:%v\n", c.RemoteAddr())
				return &TcpConn{
					o:           o,
					c:           c.(*net.TCPConn),
					listener:    nil,
					headSendBuf: make([]byte, 4),
					headRecvBuf: make([]byte, 4)}, nil
			}
		}
	}
	return nil, errors.New("new TcpConn error")
}
func (t *TcpConn) SendData(key string, val []byte) (int, error) {
	n := len(val)
	if n == 0 {
		return 0, nil
	}
	//head := make([]byte, 4)
	binary.LittleEndian.PutUint32(t.headSendBuf, uint32(n))
	n, err := t.c.Write(t.headSendBuf)
	if err != nil {
		return n, err
	}
	return t.c.Write(val)
}

func (t *TcpConn) RecvData(key string) ([]byte, error) {
	//head := make([]byte, 4)
	//todo:读取超时控制
	//n, err := t.c.Read(t.headRecvBuf)
	//使用io.readfull方法不需要特殊处理
	_, err := io.ReadFull(t.c, t.headRecvBuf)
	//if err != nil || n != 4 {
	//	return nil, errors.New("read head error")
	//}
	if err != nil {
		return nil, err
	}
	nLen := binary.LittleEndian.Uint32(t.headRecvBuf)
	data := make([]byte, nLen)
	//n, err = t.c.Read(data)
	_, err = io.ReadFull(t.c, data)
	//if err != nil || uint32(n) != nLen {
	//	return nil, errors.New("//read data buffer error!")
	//}
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (t *TcpConn) Close() error {
	if t.c != nil {
		if err := t.c.Close(); err != nil {
			return err
		}
	}
	if t.listener != nil {
		if err := t.listener.Close(); err != nil {
			return err
		}
	}
	return nil
}
