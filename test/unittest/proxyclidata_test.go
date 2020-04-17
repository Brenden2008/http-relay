package unittest

import (
	"bytes"
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func newProxyCliData() (proxyCliData *model.ProxyCliData, r *http.Request, data []byte, path string) {
	data = bytes.Repeat([]byte{10}, 10000)
	path = "/123/test"
	r, _ = http.NewRequest(http.MethodPost, "https://domain/proxy/123/test", bytes.NewReader(data))
	proxyCliData = model.NewProxyCliData(r, path)
	return
}

func TestNewProxyCliData(t *testing.T) {
	pcd, r, data, path := newProxyCliData()

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
		if bytes.Compare(b, data) != 0 {
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
