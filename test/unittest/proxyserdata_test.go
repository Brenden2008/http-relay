package unittest

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func newProxySerData() (proxyCliData *model.ProxySerData, r *http.Request, dataStr string) {
	dataStr = newString("A", 10000)
	r, _ = http.NewRequest(http.MethodPost, "https://domain/proxy/123/test", strings.NewReader(dataStr))
	proxyCliData = model.NewProxySerData(r)
	return
}

func TestNewProxySerData(t *testing.T) {
	psd, r, dataStr := newProxySerData()

	if psd.Header != &r.Header {
		t.Fail()
	}

	if b, err := ioutil.ReadAll(psd.Body); err == nil {
		if string(b) != dataStr {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
