package model

import (
	"gitlab.com/jonas.jasas/buffreader"
	"net/http"
)

type ProxyData struct {
	Path     string
	Header   *http.Header
	Body     *buffreader.BuffReader
	RespChan chan *ProxyData
}

func NewProxyData(r *http.Request, path string) *ProxyData {
	br := buffreader.New(r.Body)
	br.Buff()
	return &ProxyData{
		Path:     path,
		Header:   &r.Header,
		Body:     br,
		RespChan: make(chan *ProxyData),
	}
}
