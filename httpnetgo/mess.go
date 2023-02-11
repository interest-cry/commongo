package httpnetgo

type NetReqInfo struct {
	Data []byte `json:"data"`
}
type NetRspInfo struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	//SessId  string `json:"sess_id"`
}
type MemoryInfo struct {
	ch         chan *NetReqInfo
	remoteIp   string
	remoteport string
}
