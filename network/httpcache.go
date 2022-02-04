package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

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

//http conn req
type HttpConnRequest struct {
	Key  string `json:"key"`
	Data []byte `json:"data"`
}
type HttpConn struct {
	o          *Options
	hBigC      *HttpBigCache
	httpClient *http.Client
	tick       *time.Ticker
}

func newHttpConn(opts ...Option) (*HttpConn, error) {
	o := newOptions(opts...)
	return &HttpConn{
		o:          o,
		hBigC:      o.HttpBigC,
		httpClient: &http.Client{Transport: http.DefaultTransport},
		tick:       time.NewTicker(time.Second * time.Duration(o.TimeOut))}, nil
}

func (h *HttpConn) SendData(key string, val []byte) (int, error) {
	req := HttpConnRequest{
		Key:  key,
		Data: val,
	}
	dataJson, _ := json.Marshal(&req)
	rsp, err := h.httpClient.Post(h.o.SendUrl, "application/json", bytes.NewReader(dataJson))
	if err != nil {
		return 0, err
	}
	rsp.Body.Close()
	return len(val), nil
}
func (h *HttpConn) RecvData(key string) ([]byte, error) {
	for {
		select {
		case <-h.tick.C:
			return nil, errors.New("RecvData timeout")
		default:
			data, err := h.hBigC.bigC.Get(key)
			if err == nil {
				//todo
				//h.bigC.Delete(key)
				return data, nil
			} else if err != bigcache.ErrEntryNotFound {
				return nil, err
			}
		}
	}
}
func (h *HttpConn) Close() error {
	h.tick.Stop()
	return nil
}

//封装BigCache
type HttpBigCache struct {
	bigC *bigcache.BigCache
}

func NewHttpBigCache(sec int) *HttpBigCache {
	conf := bigcache.DefaultConfig(time.Duration(sec) * time.Second)
	conf.CleanWindow = time.Millisecond * 500
	DeLog.Infof(INFOPREFIX+"NewHttpBigCache config:%+v", conf)
	bigC, err := bigcache.NewBigCache(conf)
	if err != nil {
		panic(err)
	}
	return &HttpBigCache{
		bigC: bigC}
}

var DefaultHttpBigCache *HttpBigCache = &HttpBigCache{}

func init() {
	conf := bigcache.DefaultConfig(1800 * time.Second)
	conf.CleanWindow = time.Millisecond * 500
	DeLog.Infof(INFOPREFIX+"DefaultBigCache config:%+v", conf)
	var err error
	DefaultHttpBigCache.bigC, err = bigcache.NewBigCache(conf)
	if err != nil {
		panic(err)
	}
}

func (hb *HttpBigCache) BigCacheHandlerFunc(c *gin.Context) {
	var req HttpConnRequest
	err := c.BindJSON(&req)
	if err != nil {
		DeLog.Infof(INFOPREFIX+"SaveData BindJson error:%v", err)
		return
	}
	err = hb.bigC.Set(req.Key, req.Data)
	if err != nil {
		DeLog.Infof(INFOPREFIX+"SaveData set val error:%v", err)
		return
	}
	//DeLog.Infof(INFOPREFIX + "save data ok")
	return
	//message.Log.Infof("===>>本地缓存 set data ok")
	//todo:不需要响应返回
	//c.JSON(200, gin.H{
	//	"msg": "ok",
	//})
}
