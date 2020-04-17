package integration

import (
	"bytes"
	"fmt"
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/test/testlib"
	"gitlab.com/jonas.jasas/rwmock"
	"math/rand"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type proxyTestData [][]byte

func TestProxy(t *testing.T) {
	proxyData := genProxyData()
	ctrl := testlib.NewProxyCtrl()

}

func genProxyData() (data proxyTestData) {
	servers := []string{"first", "second", "third", "fourth", "fifth"}
	data = make(proxyTestData, 100)
	for i := 0; i < len(data); i++ {
		b := make([]byte, rand.Intn(1000000))
		rand.Read(b)
		data[i] = b
	}
	return
}

func runProxyCliReq(proxyCtrl *controller.ProxyCtrl, servers []string, data proxyTestData, closeChan chan struct{}) {
	for i := 0; i < 100; i++ {
		for _, ser := range servers {
			w, idx, path := newProxyCliReq(proxyCtrl, data, ser)
		}
	}
}

func newProxyCliReq(proxyCtrl *controller.ProxyCtrl, data proxyTestData, ser string) (w *httptest.ResponseRecorder, dataIdx int, path string) {
	path = genId(20)
	dataIdx = rand.Intn(len(data) - 1)
	url := fmt.Sprintf("https://example.com/proxy/%s/%s", ser, path)
	header := map[string]string{
		"test-path":     path,
		"test-data-idx": strconv.Itoa(dataIdx),
	}

	r := rwmock.NewShaperRand(bytes.NewReader(data[dataIdx]), 1, len(data[dataIdx]), 0, time.Second)
	w = testlib.ProxyCtrlCliReq(proxyCtrl, url, header, r)
	return
}

func newProxySerReq(proxyCtrl *controller.ProxyCtrl, data proxyTestData, ser, path, jobId string) (w *httptest.ResponseRecorder, respJobId string) {
	dataIdx := rand.Intn(len(data) - 1)
	url := fmt.Sprintf("https://example.com/proxy/%s", ser)
	header := map[string]string{
		"httprelay-proxy-headers": "test-data-idx, test-path",
		"test-path":               path,
		"test-data-idx":           strconv.Itoa(dataIdx),
	}

	r := rwmock.NewShaperRand(bytes.NewReader(data[dataIdx]), 1, len(data[dataIdx]), 0, time.Second)
	w, respJobId = testlib.ProxyCtrlSerReq(proxyCtrl, url, header, r, jobId, "")
	return
}
