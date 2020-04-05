package model

import (
	"gitlab.com/jonas.jasas/buffreader"
	"net/http"
)

type ProxyCliData struct {
	Method   string
	Path     string
	Header   *http.Header
	Body     *buffreader.BuffReader
	RespChan chan *ProxySerData
}

func NewProxyCliData(r *http.Request, path string) (proxyReqData *ProxyCliData) {
	proxyReqData = &ProxyCliData{
		Method: r.Method,
		Path:   path,
		Header: &r.Header,
		Body:   buffreader.New(r.Body),
	}

	proxyReqData.Body.Buff()
	proxyReqData.RespChan = make(chan *ProxySerData)

	return
}

type ProxySerData struct {
	Header *http.Header
	Body   *buffreader.BuffReader
}

func NewProxySerData(r *http.Request) (proxyRespData *ProxySerData) {
	proxyRespData = &ProxySerData{
		Header: &r.Header,
		Body:   buffreader.New(r.Body),
	}
	proxyRespData.Body.Buff()
	return
}
