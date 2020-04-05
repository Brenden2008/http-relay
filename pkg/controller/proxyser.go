package controller

import (
	"errors"
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"io"
	"net/http"
)

func (pc *ProxyCtrl) handleServer(ser *model.ProxySer, r *http.Request, w http.ResponseWriter) {
	if jobId := r.Header.Get("Httprelay-Proxy-Jobid"); jobId != "" {
		if cliData, ok := ser.TakeJob(jobId); ok { // Request is previous job response /////////////////////////////////////
			serData := model.NewProxySerData(r)
			if pc.transferSerReq(cliData.RespChan, serData, r, w) != nil {
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

func (pc *ProxyCtrl) transferSerReq(respChan chan<- *model.ProxySerData, data *model.ProxySerData, r *http.Request, w http.ResponseWriter) (err error) {
	select {
	case respChan <- data:
		close(respChan)
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

func (pc *ProxyCtrl) transferSerResp(ser *model.ProxySer, r *http.Request, w http.ResponseWriter) {
	select { // Response is new job request ////////////////////////////////////////////////////////////////////////////
	case cliData := <-ser.ReqChan:
		jobId := randStr(8)
		w.Header().Add("Httprelay-Proxy-Jobid", jobId)
		w.Header().Add("Httprelay-Proxy-Method", cliData.Method)
		w.Header().Add("Httprelay-Proxy-Path", cliData.Path)

		if err := cliData.Header.Write(w); err == nil {
			if _, err := io.Copy(w, cliData.Body); err == nil {
				ser.AddJob(jobId, cliData)
			} else {
				//TODO: Log body transfer err
				return
			}
		} else {
			//TODO: Log header transfer err
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
