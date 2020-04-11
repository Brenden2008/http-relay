package unittest

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func newProxyCliData() (proxyCliData *model.ProxyCliData, r *http.Request, dataStr string, path string) {
	dataStr = newString("A", 10000)
	path = "/123/test"
	r, _ = http.NewRequest(http.MethodPost, "https://domain/proxy/123/test", strings.NewReader(dataStr))
	proxyCliData = model.NewProxyCliData(r, path)
	return
}

func TestNewProxyCliData(t *testing.T) {
	pcd, r, dataStr, path := newProxyCliData()

	if pcd.Method != http.MethodPost {
		t.Fail()
	}

	if pcd.Path != path {
		t.Fail()
	}

	if pcd.Header != &r.Header {
		t.Fail()
	}

	if b, err := ioutil.ReadAll(pcd.Body); err == nil {
		if string(b) != dataStr {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}

func TestProxyCliDataClose(t *testing.T) {
	pcd, _, _, _ := newProxyCliData()
	if !pcd.CloseRespChan() {
		t.Fail()
	}
	if pcd.CloseRespChan() {
		t.Fail()
	}

	pcd, _, _, _ = newProxyCliData()
	go func() { pcd.RespChan <- nil }()
	time.Sleep(100000)
	if !pcd.CloseRespChan() {
		t.Fail()
	}
	if pcd.CloseRespChan() {
		t.Fail()
	}
}
