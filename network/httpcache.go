package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

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

func newHttpConn(o *Options) (*HttpConn, error) {
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
			return nil, errors.New("recv data timeout")
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

//封装bigcache
type HttpBigCache struct {
	bigC *bigcache.BigCache
}

func NewHttpBigCache(sec int) *HttpBigCache {
	conf := bigcache.DefaultConfig(time.Duration(sec) * time.Second)
	conf.CleanWindow = time.Millisecond * 500
	fmt.Printf("...NewHttpBigCache config:%+v\n", conf)
	bigC, err := bigcache.NewBigCache(conf)
	if err != nil {
		panic(err)
	}
	return &HttpBigCache{
		bigC: bigC,
	}
}

var DefaultHttpBigCache *HttpBigCache = &HttpBigCache{}

func init() {
	conf := bigcache.DefaultConfig(1800 * time.Second)
	conf.CleanWindow = time.Millisecond * 500
	fmt.Printf("DefaultBigCache config:%+v\n", conf)
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
		fmt.Printf("SaveData BindJson error:%v\n", err)
		return
	}
	err = hb.bigC.Set(req.Key, req.Data)
	if err != nil {
		fmt.Printf("SaveData set val error:%v", err)
		return
	}
	fmt.Printf("save data ok\n")
	//message.Log.Infof("===>>本地缓存 set data ok")
	//todo:不需要响应返回
	//c.JSON(200, gin.H{
	//	"msg": "ok",
	//})
}
