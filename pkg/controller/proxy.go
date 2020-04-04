package controller

import (
	"errors"
	"fmt"
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type ProxyCtrl struct {
	rep      repository.ProxyRep
	stopChan <-chan struct{}
	*model.Waiters
}

func NewProxyCtrl(rep repository.ProxyRep, stopChan <-chan struct{}) *ProxyCtrl {
	return &ProxyCtrl{
		rep:      rep,
		stopChan: stopChan,
		Waiters:  model.NewWaiters(),
	}
}

func (pc *ProxyCtrl) Conduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	select {
	case <-pc.stopChan:
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	default:
	}

	pathArr := strings.Split(r.URL.Path, "/")
	ser := pc.rep.GetSer(pathArr[1])

	if jobId, ok := jobId(r); ok { // Server //////////////////////////////////
		if cliData, ok := ser.TakeJob(jobId); ok { // Has response from previous job
			serData := model.NewProxyData(r, "")
			if pc.transferReq(cliData.RespChan, serData, r.Context().Done()) != nil {
				//TODO: Log err
			}
		} else {
			//TODO: Log job not found
		}

		select {
		case cliData := <-ser.ReqChan:
			if pc.transferResp(cliData, w, r.Context().Done()) != nil {
				//TODO: Log err
			}
		case <-pc.stopChan:
			w.WriteHeader(http.StatusServiceUnavailable)
			err = errors.New("Proxy controller transferReq stop signal received")
		case <-closeChan:
			w.WriteHeader(http.StatusServiceUnavailable)
			err = errors.New("Proxy controller transferReq close signal received")
		}

	} else { // Client //////////////////////////////////
		pathPrefix := fmt.Sprintf("/%s/%s", pathArr[0], pathArr[1])
		serReqPath := strings.TrimPrefix(r.URL.Path, pathPrefix)

		data := model.NewProxyData(r, serReqPath)
		defer close(data.RespChan)

		if pc.transferResp(data, w, r.Context().Done()) != nil {
			//TODO: Log err
		}
	}
}

func (pc *ProxyCtrl) transferReq(dataChan chan<- *model.ProxyData, data *model.ProxyData, closeChan <-chan struct{}) (err error) {
	select {
	case dataChan <- data:
	case <-pc.stopChan:
		err = errors.New("Proxy controller transferReq stop signal received")
	case <-closeChan:
		err = errors.New("Proxy controller transferReq close signal received")
	}
}

func (pc *ProxyCtrl) transferResp(data *model.ProxyData, w http.ResponseWriter, closeChan <-chan struct{}) (err error) {
	select {
	case respData := <-data.RespChan:
		if err = respData.Header.Write(w); err == nil {
			_, err = io.Copy(w, respData.Body)
		}
	case <-pc.stopChan:
		err = errors.New("Proxy controller transferResp stop signal received")
	case <-closeChan:
		err = errors.New("Proxy controller transferResp close signal received")
	}
}

func jobId(r *http.Request) (jobId string, ok bool) {
	vals := r.URL.Query()["jobid"]
	if len(vals) > 0 {
		jobId = vals[0]
	}
	return
}

// https://stackoverflow.com/a/22892986/625521 ////////////////////////////////////////
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//////////////////////////////////////////////////////////////////////////////////////
