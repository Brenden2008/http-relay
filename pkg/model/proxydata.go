package model

import (
	"gitlab.com/jonas.jasas/buffreader"
	"net/http"
	"sync"
)

type ProxyCliData struct {
	Method    string
	Path      string
	Header    *http.Header
	Body      *buffreader.BuffReader
	RespChan  chan *ProxySerData
	RespChanL sync.Mutex
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

func (pcd *ProxyCliData) CloseRespChan() (ok bool) {
	pcd.RespChanL.Lock()
	defer pcd.RespChanL.Unlock()

	select {
	case _, ok = <-pcd.RespChan:
	default:
		ok = true
	}

	if ok {
		close(pcd.RespChan)
	}

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
