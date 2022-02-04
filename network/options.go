package network

type Options struct {
	//NetWorkType    string
	Ip             string
	Port           int
	TimeOut        int
	ClientOrServer string
	//httpcache
	HttpBigC *HttpBigCache
	SendUrl  string
	//EventBus
	EventB *EventBus
}

func newOptions(opts ...Option) *Options {
	opt := Options{
		//NetWorkType:    TCP,
		Ip:             "127.0.0.1",
		Port:           18888,
		TimeOut:        3600,
		ClientOrServer: CLIENT,
		HttpBigC:       DefaultHttpBigCache,
		EventB:         DefaultEventBus,
	}
	for _, o := range opts {
		o(&opt)
	}
	DeLog.Infof(INFOPREFIX+"Options:%+v\n", opt)
	return &opt
}

//func NetWorkType(netType string) Option {
//	return func(o *Options) {
//		o.NetWorkType = netType
//	}
//}
func TimeOut(seconds int) Option {
	return func(o *Options) {
		o.TimeOut = seconds
	}
}
