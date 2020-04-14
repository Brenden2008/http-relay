package unittest

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var cliData1 = newString("Client data 1. ", 10000)
var cliData2 = newString("Client data 2. ", 10000)
var serData1 = newString("Server data 1. ", 10000)
var serData2 = newString("Server data 2. ", 10000)

const serUrl = "https://domain/proxy/123/test"

func newReq(method string, url string, data string) (r *http.Request) {
	r, _ = http.NewRequest(method, url, strings.NewReader(data))
	return
}

func respDataEq(body io.Reader, data string) bool {
	respData, _ := ioutil.ReadAll(body)
	return string(respData) == data
}

func newCtrl() (proxyCtrl *controller.ProxyCtrl, stopChan chan struct{}, closeChan chan struct{}) {
	stopChan = make(chan struct{})
	proxyRep := repository.NewProxyRep()
	proxyCtrl = controller.NewProxyCtrl(proxyRep, stopChan)
	closeChan = make(chan struct{})
	return
}

func proxyCtrlCliReq(ctrl *controller.ProxyCtrl, data string) *httptest.ResponseRecorder {
	r := newReq(http.MethodPost, serUrl, data)
	w := httptest.NewRecorder()
	ctrl.Conduct(w, r)
	return w
}
func proxyCtrlSerReq(ctrl *controller.ProxyCtrl, data string, reqJobId string, wSecret string) (resp *httptest.ResponseRecorder, respJobId string) {
	r := newReq("SERVE", serUrl, data)
	r.Header.Add("httprelay-proxy-jobid", reqJobId)
	r.Header.Add("httprelay-wsecret", wSecret)
	w := httptest.NewRecorder()
	ctrl.Conduct(w, r)
	respJobId = w.Header().Get("httprelay-proxy-jobid")
	return w, respJobId
}

/////////////////////////////////////////////////////////////////////////

func TestProxyCtrlConduct(t *testing.T) {
	proxyCtrl, _, closeChan := newCtrl()
	go func() {
		defer close(closeChan)
		resp := proxyCtrlCliReq(proxyCtrl, cliData1)
		if !respDataEq(resp.Body, serData1) {
			t.Error("Client received wrong response body")
		}
	}()
	go func() {
		defer close(closeChan)
		resp, jobId := proxyCtrlSerReq(proxyCtrl, "", "", "")
		if !respDataEq(resp.Body, cliData1) {
			t.Error("Server received wrong client data body")
			return
		}
		proxyCtrlSerReq(proxyCtrl, serData1, jobId, "")
	}()
	<-closeChan
}

func TestProxyCtrlWSecret(t *testing.T) {
	proxyCtrl, _, closeChan := newCtrl()
	go func() {
		defer close(closeChan)
		resp := proxyCtrlCliReq(proxyCtrl, cliData1)
		if !respDataEq(resp.Body, serData1) {
			t.Error("Client received wrong response body")
			return
		}
		resp = proxyCtrlCliReq(proxyCtrl, cliData2)
		if !respDataEq(resp.Body, serData2) {
			t.Error("Client received wrong response body")
			return
		}
	}()

	go func() {
		const goodSecret = "secret1"
		const badSecret = "bad secret"
		defer close(closeChan)
		resp, jobId := proxyCtrlSerReq(proxyCtrl, "", "", goodSecret)
		if !respDataEq(resp.Body, cliData1) {
			t.Error("Server received wrong client data body")
			return
		}

		resp, _ = proxyCtrlSerReq(proxyCtrl, serData1, jobId, badSecret)
		if resp.Code != http.StatusUnauthorized {
			t.Error("Server is accessing unauthorized data channel with the bad secret")
			return
		}

		resp, jobId = proxyCtrlSerReq(proxyCtrl, serData1, jobId, goodSecret)
		if !respDataEq(resp.Body, cliData2) {
			t.Error("Server received wrong client data body")
			return
		}

		resp, _ = proxyCtrlSerReq(proxyCtrl, serData1, jobId, badSecret)
		if resp.Code != http.StatusUnauthorized {
			t.Error("Server is accessing unauthorized data channel with the bad secret")
			return
		}

		resp, jobId = proxyCtrlSerReq(proxyCtrl, serData2, jobId, goodSecret)
	}()
	<-closeChan
}
