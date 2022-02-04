package server

//func TestServer_SaveData(t *testing.T) {
//	var relativePath string = "/v1/sendppp"
//	ser, err := NewServer(":8900")
//	assert.NoError(t, err)
//	route := func() *gin.Engine {
//		r := gin.Default()
//		r.POST(relativePath, func(c *gin.Context) {
//			ser.SaveData(c)
//		})
//		return r
//	}()
//	serTest := httptest.NewServer(route)
//	//serTest := httptest.NewServer(http.HandlerFunc(route.ServeHTTP))
//	defer serTest.Close()
//	clent := http.Client{}
//	req := network.HttpCacheRequest{Key: "key", Data: []byte("ddddddddddd")}
//	jq, err := json.Marshal(&req)
//	assert.NoError(t, err)
//	rsp, err := clent.Post(serTest.URL+relativePath, "application/json", bytes.NewReader(jq))
//	assert.NoError(t, err)
//	defer func() { rsp.Body.Close() }()
//	out, err := ioutil.ReadAll(rsp.Body)
//	assert.NoError(t, err)
//	fmt.Printf("===>>out:%v\n", string(out))
//}
