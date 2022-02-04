package network

import (
	"fmt"
)

type Options struct {
	NetWorkType    string
	Ip             string
	Port           int
	TimeOut        int
	ClientOrServer string
	//httpcache
	HttpBigC *HttpBigCache
	SendUrl  string
}

func newOptions(opts ...Option) *Options {
	opt := Options{
		NetWorkType:    TCP,
		Ip:             "127.0.0.1",
		Port:           18888,
		TimeOut:        3600,
		ClientOrServer: CLIENT,
		HttpBigC:       DefaultHttpBigCache,
	}
	for _, o := range opts {
		o(&opt)
	}
	//if opt.NetWorkType == TCP {
	//	//if opt.ClientOrServer != CLIENT && opt.ClientOrServer != SERVER {
	//	//	panic("tcp network must set client or server")
	//	//}
	//}
	fmt.Printf("***options:%+v\n", opt)
	return &opt
}
func NetWorkType(netType string) Option {
	return func(o *Options) {
		o.NetWorkType = netType
	}
}
func Ip(ip string) Option {
	return func(o *Options) {
		o.Ip = ip
	}
}
func Port(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}
func TimeOut(seconds int) Option {
	return func(o *Options) {
		o.TimeOut = seconds
	}
}
func ClientOrServer(cliOrSer string) Option {
	return func(o *Options) {
		o.ClientOrServer = cliOrSer
	}
}

//httpcache
func BigCache(bigCache *HttpBigCache) Option {
	return func(o *Options) {
		o.HttpBigC = bigCache
	}
}
func SendUrl(sendUrl string) Option {
	return func(o *Options) {
		o.SendUrl = sendUrl
	}
}
