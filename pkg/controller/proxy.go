package controller

import (
	"errors"
	"fmt"
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"io"
	"net/http"
	"strings"
)

type ProxyCtrl struct {
	rep      *repository.ProxyRep
	stopChan <-chan struct{}
	*model.Waiters
}

func NewProxyCtrl(rep *repository.ProxyRep, stopChan <-chan struct{}) *ProxyCtrl {
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
	ser := pc.rep.GetSer(pathArr[2])

	if strings.EqualFold(r.Method, "SERVE") {
		pc.handleServer(ser, r, w)
	} else {
		pc.handleClient(ser, pathArr, r, w)
	}
}

func (pc *ProxyCtrl) handleClient(ser *model.ProxySer, pathArr []string, r *http.Request, w http.ResponseWriter) {
	pathPrefix := fmt.Sprintf("/%s/%s", pathArr[1], pathArr[2])
	serReqPath := strings.TrimPrefix(r.URL.Path, pathPrefix)
	if serReqPath == "" {
		serReqPath = "/"
	}

	data := model.NewProxyData(r, serReqPath)

	if pc.transferReq(ser.ReqChan, data, r, w) == nil {
		if pc.transferResp(data, r, w) != nil {
			//TODO: Log err
			return
		}
	} else {
		data.Body.Close()
		//TODO: Log err
	}
}

func (pc *ProxyCtrl) handleServer(ser *model.ProxySer, r *http.Request, w http.ResponseWriter) {
	if jobId := r.Header.Get("Httprelay-Proxy-Jobid"); jobId != "" {
		if cliData, ok := ser.TakeJob(jobId); ok { // Request is previous job response /////////////////////////////////////
			serData := model.NewProxyData(r, "")
			if pc.transferReq(cliData.RespChan, serData, r, w) != nil {
				//TODO: Log request transfer err
				return
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			//TODO: Log job not found
			return
		}
	}

	select { // Response is new job request ////////////////////////////////////////////////////////////////////////////
	case cliData := <-ser.ReqChan:
		jobId := randStr(8)
		w.Header().Add("Httprelay-Proxy-Jobid", jobId)
		w.Header().Add("Httprelay-Proxy-Path", cliData.Path)
		if pc.transferResp(cliData, r, w) == nil {
			ser.AddJob(jobId, cliData)
		} else {
			//TODO: Log err
			return
		}
	case <-pc.stopChan:
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	case <-r.Context().Done():
		//TODO: log
		return
	}
}

func (pc *ProxyCtrl) transferReq(dataChan chan<- *model.ProxyData, data *model.ProxyData, r *http.Request, w http.ResponseWriter) (err error) {
	select {
	case dataChan <- data:
		close(dataChan)
	case <-pc.stopChan:
		data.Body.Close() // Stopping buffering
		w.WriteHeader(http.StatusServiceUnavailable)
		err = errors.New("proxy controller transferReq stop signal received")
	case <-r.Context().Done():
		data.Body.Close() // Stopping buffering
		w.WriteHeader(http.StatusBadGateway)
		err = errors.New("proxy controller transferReq close signal received")
	}
	return
}

func (pc *ProxyCtrl) transferResp(data *model.ProxyData, r *http.Request, w http.ResponseWriter) (err error) {
	select {
	case respData := <-data.RespChan:
		if err = respData.Header.Write(w); err == nil {
			_, err = io.Copy(w, respData.Body)
		}
	case <-pc.stopChan:
		close(data.RespChan)
		w.WriteHeader(http.StatusServiceUnavailable)
		err = errors.New("proxy controller transferResp stop signal received")
	case <-r.Context().Done():
		close(data.RespChan)
		w.WriteHeader(http.StatusBadGateway)
		err = errors.New("proxy controller transferResp close signal received")
	}
	return
}

func jobId(r *http.Request) (jobId string, ok bool) {
	vals := r.URL.Query()["jobid"]
	if len(vals) > 0 {
		ok = true
		jobId = vals[0]
	}
	return
}
