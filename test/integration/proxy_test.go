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
	"sync"
	"testing"
	"time"
)

type proxyTestData [][]byte

func TestProxy(t *testing.T) {
	servers := []string{"first", "second", "third", "fourth", "fifth"}
	proxyData := genProxyData()
	ctrl, _, _ := testlib.NewProxyCtrl()
	closeChan := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		go func() {
			wg.Add(1)
			runProxyCliReq(t, ctrl, servers, proxyData, closeChan)
			t.Log("Client done")
			wg.Done()
		}()
	}

	for _, server := range servers {
		go func(ser string) {
			runProxySerReq(t, ctrl, ser, proxyData, closeChan)
			close(closeChan)
		}(server)
	}

	wg.Wait()
	//close(stopChan)
}

func genProxyData() (data proxyTestData) {
	data = make(proxyTestData, 100)
	for i := 0; i < len(data); i++ {
		b := make([]byte, rand.Intn(1000000))
		rand.Read(b)
		data[i] = b
	}
	return
}

func runProxyCliReq(t *testing.T, proxyCtrl *controller.ProxyCtrl, servers []string, data proxyTestData, closeChan chan struct{}) {
	for i := 0; i < 10; i++ {
		for _, ser := range servers {
			w, path := newProxyCliReq(proxyCtrl, data, ser)
			select {
			case <-closeChan:
				return
			default:
			}
			if w.Header().Get("test-path") != path {
				t.Error("client received incorrect path in response")
				close(closeChan)
				return
			}
			if dataIdx, err := strconv.Atoi(w.Header().Get("test-data-idx")); err == nil {
				if testlib.RespDataEq(w.Body, data[dataIdx]) {
					t.Error("client received incorrect body in response")
					close(closeChan)
					return
				}
			} else {
				t.Error("client received incorrect data array index")
				close(closeChan)
				return
			}
		}
	}
}

func runProxySerReq(t *testing.T, proxyCtrl *controller.ProxyCtrl, server string, data proxyTestData, closeChan chan struct{}) {
	path, jobId  := "", ""
	var w *httptest.ResponseRecorder
	for {
		t.Log("Server req")
		w, jobId = newProxySerReq(proxyCtrl, data, server, path, jobId)
		t.Log("Server received req")
		path = w.Header().Get("test-path")
		jobId = w.Header().Get("httprelay-proxy-jobid")

		if dataIdx, err := strconv.Atoi(w.Header().Get("test-data-idx")); err == nil {
			if testlib.RespDataEq(w.Body, data[dataIdx]) {
				t.Error("server received incorrect body in response")
				close(closeChan)
				return
			}
		} else {
			t.Error("server received incorrect data array index")
			close(closeChan)
			return
		}

		select {
		case <-closeChan:
			return
		default:
		}
	}
}

func newProxyCliReq(proxyCtrl *controller.ProxyCtrl, data proxyTestData, ser string) (w *httptest.ResponseRecorder, path string) {
	path = genId(20)
	dataIdx := rand.Intn(len(data) - 1)
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
