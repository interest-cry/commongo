package network

type Options struct {
	Ip             string
	Port           int
	TimeOut        int
	ClientOrServer string
	HttpBigC       *HttpBigCache
	SendUrl        string
	EventB         *EventBus
	Uid            string
	LocalNid       string
	RemoteNid      string
}

func newOptions(opts ...Option) *Options {
	opt := Options{
		Ip:             "127.0.0.1",
		Port:           18888,
		TimeOut:        3600,
		ClientOrServer: CLIENT,
		HttpBigC:       DefaultHttpBigCache,
		SendUrl:        "/v1/send",
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
func ClientOrServer(cliOrSer string) Option {
	return func(o *Options) {
		o.ClientOrServer = cliOrSer
	}
}
func TimeOut(seconds int) Option {
	return func(o *Options) {
		o.TimeOut = seconds
	}
}

//Options for HttpBigCache
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

//for chanconn
func EventBusSet(eventB *EventBus) Option {
	return func(o *Options) {
		o.EventB = eventB
	}
}
func Uid(uid string) Option {
	return func(o *Options) {
		o.Uid = uid
	}
}
func LocalNid(localNid string) Option {
	return func(o *Options) {
		o.LocalNid = localNid
	}
}
func RemoteNid(remoteNid string) Option {
	return func(o *Options) {
		o.RemoteNid = remoteNid
	}
}
