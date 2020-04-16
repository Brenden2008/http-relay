package testlib

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"net/http"
	"net/http/httptest"
)

func NewProxyCtrl() (proxyCtrl *controller.ProxyCtrl, stopChan chan struct{}, closeChan chan struct{}) {
	stopChan = make(chan struct{})
	proxyRep := repository.NewProxyRep()
	proxyCtrl = controller.NewProxyCtrl(proxyRep, stopChan)
	closeChan = make(chan struct{})
	return
}

func ProxyCtrlCliReq(ctrl *controller.ProxyCtrl, url string, header map[string]string, data string) *httptest.ResponseRecorder {
	r := newReq(http.MethodPost, url, header, data)
	w := httptest.NewRecorder()
	ctrl.Conduct(w, r)
	return w
}

func ProxyCtrlSerReq(ctrl *controller.ProxyCtrl, url string, header map[string]string, data string, reqJobId string, wSecret string) (resp *httptest.ResponseRecorder, respJobId string) {
	r := newReq("SERVE", url, header, data)
	r.Header.Add("httprelay-proxy-jobid", reqJobId)
	r.Header.Add("httprelay-wsecret", wSecret)
	w := httptest.NewRecorder()
	ctrl.Conduct(w, r)
	respJobId = w.Header().Get("httprelay-proxy-jobid")
	return w, respJobId
}
