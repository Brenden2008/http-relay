package model

type ProxyData struct {
	Req      *PtpData
	RespChan chan *PtpData
}

func NewProxyData(data *PtpData) *ProxyData {
	return &ProxyData{
		Req:      data,
		RespChan: make(chan *PtpData),
	}
}
