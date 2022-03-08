package network

import (
	"github.com/allegro/bigcache"
	"github.com/sirupsen/logrus"
	"time"
)

type Options struct {
	NetworkType    string
	Ip             string
	Port           int
	TimeOut        int
	ClientOrServer string
	HttpBigcache   *HttpBigCache
	SendUrl        string
	EventB         *EventBus
	Uid            string
	LocalNid       string
	RemoteNid      string
}

func newOptions(opts ...Option) *Options {
	opt := Options{
		NetworkType:    TCPCONN,
		Ip:             "127.0.0.1",
		Port:           18888,
		TimeOut:        3600,
		ClientOrServer: CLIENT,
		HttpBigcache:   DefaultHttpBigCache,
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
func NetworkType(networkTpye string) Option {
	return func(o *Options) {
		o.NetworkType = networkTpye
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
		o.HttpBigcache = bigCache
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

var (
	DeLog *logrus.Logger = logrus.New()
	Log   *logrus.Logger = logrus.New()
)

const (
	INFOPREFIX = "[=== INFO]"
	WARNPREFIX = "[=== WARN]"
)

func init() {
	DeLog.SetLevel(logrus.DebugLevel)
	Log.SetLevel(logrus.ErrorLevel)
	conf := bigcache.DefaultConfig(1800 * time.Second)
	conf.CleanWindow = time.Millisecond * 500
	DeLog.Infof(INFOPREFIX+"DefaultBigCache config:%+v", conf)
	var err error
	DefaultHttpBigCache.bigCache, err = bigcache.NewBigCache(conf)
	if err != nil {
		panic(err)
	}
	DefaultGinSenderMap.ginSenders.Store(CACHECONN, DefaultHttpBigCache)
	//DefaultGinSenderMap.ginSenders[CACHECONN] = DefaultHttpBigCache
	DefaultGinSenderMap.ginSenders.Store(CHANCONN, DefaultEventBus)
	//DefaultGinSenderMap.ginSenders[CHANCONN] = DefaultEventBus
}
