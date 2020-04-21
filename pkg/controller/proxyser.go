package controller

import (
	"errors"
	"fmt"
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"io"
	"net/http"
)

func (pc *ProxyCtrl) handleServer(ser *model.ProxySer, r *http.Request, w http.ResponseWriter) {
	if !ser.WAuth(wSecret(r)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if jobId := r.Header.Get("Httprelay-Proxy-Jobid"); jobId != "" {
		if cliData, ok := ser.TakeJob(jobId); ok { // Request is previous job response /////////////////////////////////////
			defer cliData.CloseRespChan()

			serData := model.NewProxySerData(r)
			if pc.transferSerReq(cliData.RespChan, serData, r, w) != nil {
				if serData.Body.Close() != nil { // Stopping buffering
					//TODO: Log buffering error (to silence compiler warnings)
				}
				//TODO: Log request transfer err
				return
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			//TODO: Log job not found
			return
		}
	}

	pc.transferSerResp(ser, r, w)
}

func (pc *ProxyCtrl) transferSerReq(respChan chan<- *model.ProxySerData, serData *model.ProxySerData, r *http.Request, w http.ResponseWriter) (err error) {
	select {
	case respChan <- serData:
	case <-pc.stopChan:
		w.WriteHeader(http.StatusServiceUnavailable)
		err = errors.New("proxy controller transferReq stop signal received")
	case <-r.Context().Done():
		w.WriteHeader(http.StatusBadGateway)
		err = errors.New("proxy controller transferReq close signal received")
	}
	return
}

func (pc *ProxyCtrl) transferSerResp(ser *model.ProxySer, r *http.Request, w http.ResponseWriter) {
	select { // Response is new job request ////////////////////////////////////////////////////////////////////////////
	case cliData := <-ser.ReqChan:
		jobId := randStr(8)

		for name, vals := range *cliData.Header {
			for _, val := range vals {
				w.Header().Add(name, val)
			}
		}

		w.Header().Set("Httprelay-Proxy-Jobid", jobId)
		w.Header().Set("Httprelay-Proxy-Method", cliData.Method)
		w.Header().Set("Httprelay-Proxy-Path", cliData.Path)

		if _, err := io.Copy(w, cliData.Body); err == nil {
			ser.AddJob(jobId, cliData)
		} else {
			//TODO: Log body transfer err
			return
		}

	case <-pc.stopChan:
		fmt.Println("stop in transferSerResp")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	case <-r.Context().Done():
		fmt.Println("close in transferSerResp")
		//TODO: log
		return
	}
}
