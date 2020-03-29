package model

import "net/http"

type ProxyData struct {
	CliChan chan *PtpData
	SerChan chan *PtpData
}

func NewProxyData(r *http.Request) *ProxyData {
	return &ProxyData{
		CliChan: make(chan *PtpData),
		SerChan: make(chan *PtpData),
	}
}
