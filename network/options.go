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
	EventB    *EventBus
	Uid       string
	LocalNid  string
	RemoteNid string
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
		Uid:            "2022-02-04",
		LocalNid:       "local-node-001",
		RemoteNid:      "remote-node-001",
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
