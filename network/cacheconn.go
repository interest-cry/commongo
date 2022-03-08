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

type CacheConn struct {
	o            *Options
	httpBigCache *HttpBigCache
	httpClient   *http.Client
	tick         *time.Timer
	timeout      time.Duration
}

func newCacheConn(opts ...Option) (Communicator, error) {
	o := newOptions(opts...)
	timeout := time.Second * time.Duration(o.TimeOut)
	return &CacheConn{
		o:            o,
		httpBigCache: o.HttpBigcache,
		httpClient:   &http.Client{Transport: http.DefaultTransport},
		tick:         time.NewTimer(timeout),
		timeout:      timeout}, nil
}

func (h *CacheConn) SendData(key string, val []byte) (int, error) {
	req := requestForSend{
		NetworkType: h.o.NetworkType,
		Key:         key,
		Data:        val,
	}
	dataJson, _ := json.Marshal(&req)
	rsp, err := h.httpClient.Post(h.o.SendUrl, "application/json", bytes.NewReader(dataJson))
	if err != nil {
		return 0, err
	}
	rsp.Body.Close()
	return len(val), nil
}
func (h *CacheConn) RecvData(key string) ([]byte, error) {
	defer func() {
		h.tick.Reset(h.timeout)
	}()
	for {
		select {
		case <-h.tick.C:
			return nil, errors.New("RecvData timeout")
		default:
			data, err := h.httpBigCache.bigCache.Get(key)
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
func (h *CacheConn) Close() error {
	h.tick.Stop()
	return nil
}

//封装BigCache
type HttpBigCache struct {
	bigCache *bigcache.BigCache
}

func NewHttpBigCache(sec int) *HttpBigCache {
	conf := bigcache.DefaultConfig(time.Duration(sec) * time.Second)
	conf.CleanWindow = time.Millisecond * 500
	DeLog.Infof(INFOPREFIX+"NewHttpBigCache config:%+v", conf)
	bigCache, err := bigcache.NewBigCache(conf)
	if err != nil {
		panic(err)
	}
	return &HttpBigCache{
		bigCache: bigCache}
}

var DefaultHttpBigCache *HttpBigCache = &HttpBigCache{}

func (hb *HttpBigCache) WriteMessageGin(c *gin.Context) {
	var req requestForSend
	err := c.BindJSON(&req)
	if err != nil {
		DeLog.Infof(INFOPREFIX+"SaveData BindJson error:%v", err)
		return
	}
	err = hb.bigCache.Set(req.Key, req.Data)
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

func (hb *HttpBigCache) WriteMessage(req requestForSend) {
	err := hb.bigCache.Set(req.Key, req.Data)
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
