package unittest

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var cliData = newString("Client data. ", 10000)
var serData = newString("Server data. ", 10000)

const serUrl = "https://domain/proxy/123/test"

func newReq(method string, url string, data string) (r *http.Request) {
	r, _ = http.NewRequest(method, url, strings.NewReader(data))
	return
}

func TestProxyCtrlConduct(t *testing.T) {
	stopChan := make(chan struct{})
	proxyRep := repository.NewProxyRep()
	proxyCtrl := controller.NewProxyCtrl(proxyRep, stopChan)

	closeChan := make(chan struct{})
	go func() {
		proxyCtrlCliReq(t, proxyCtrl, closeChan)
	}()

	go func() {
		proxyCtrlSerReqResp(t, proxyCtrl, closeChan)
	}()

	<-closeChan
}

func proxyCtrlCliReq(t *testing.T, ctrl *controller.ProxyCtrl, closeChan chan struct{}) {
	r := newReq(http.MethodPost, serUrl, cliData)
	w := httptest.NewRecorder()
	ctrl.Conduct(w, r)
	body, _ := ioutil.ReadAll(w.Body)
	if string(body) != serData {
		t.Error("Client received wrong response body")
	}
	close(closeChan)
}

func proxyCtrlSerReqResp(t *testing.T, ctrl *controller.ProxyCtrl, closeChan chan struct{}) {
	r1 := newReq("SERVE", serUrl, "")
	w1 := httptest.NewRecorder()
	ctrl.Conduct(w1, r1)
	body, _ := ioutil.ReadAll(w1.Body)
	if string(body) != cliData {
		t.Error("Server received wrong client data body")
		close(closeChan)
		return
	}

	r2 := newReq("SERVE", serUrl, serData)
	r2.Header.Add("httprelay-proxy-jobid", w1.Header().Get("httprelay-proxy-jobid"))
	w2 := httptest.NewRecorder()
	ctrl.Conduct(w2, r2)
}
