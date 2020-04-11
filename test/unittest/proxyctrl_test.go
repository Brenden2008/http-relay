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

func newReq(method string, url string, data string) (r *http.Request, dataStr string) {
	dataStr = newString(data, 10000)
	r, _ = http.NewRequest(method, url, strings.NewReader(dataStr))
	return
}

func TestProxyCtrlConduct(t *testing.T) {
	stopChan := make(chan struct{})
	proxyRep := repository.NewProxyRep()
	proxyCtrl := controller.NewProxyCtrl(proxyRep, stopChan)

	cliCloseChan := make(chan struct{})
	go func() {
		proxyCtrlCliReq(t, proxyCtrl)
		close(cliCloseChan)
	}()

	serCloseChan := make(chan struct{})
	go func() {
		proxyCtrlSerReq(t, proxyCtrl)
		proxyCtrlSerResp(t, proxyCtrl)
		close(serCloseChan)
	}()

	<-serCloseChan
	<-cliCloseChan
}

func proxyCtrlCliReq(t *testing.T, ctrl *controller.ProxyCtrl) {
	r, _ := newReq(http.MethodPost, serUrl, cliData)
	w := httptest.NewRecorder()
	ctrl.Conduct(w, r)
	body, _ := ioutil.ReadAll(w.Body)
	if string(body) != serData {
		t.Fail()
	}
}

func proxyCtrlSerReq(t *testing.T, ctrl *controller.ProxyCtrl) {
	r, _ := newReq("SERVE", serUrl, cliData)
	w := httptest.NewRecorder()
	ctrl.Conduct(w, r)
	body, _ := ioutil.ReadAll(w.Body)
	if string(body) != serData {
		t.Fail()
	}
}

func proxyCtrlSerResp(t *testing.T, ctrl *controller.ProxyCtrl) {
	r, _ := newReq("SERVE", serUrl, cliData)
	w := httptest.NewRecorder()
	ctrl.Conduct(w, r)
	body, _ := ioutil.ReadAll(w.Body)
	if string(body) != serData {
		t.Fail()
	}
}
