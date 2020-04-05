package controller

import (
	"errors"
	"fmt"
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"io"
	"net/http"
	"strings"
)

func (pc *ProxyCtrl) handleClient(ser *model.ProxySer, pathArr []string, r *http.Request, w http.ResponseWriter) {
	pathPrefix := fmt.Sprintf("/%s/%s", pathArr[1], pathArr[2])
	serReqPath := strings.TrimPrefix(r.URL.Path, pathPrefix)
	if serReqPath == "" {
		serReqPath = "/"
	}

	cliData := model.NewProxyCliData(r, serReqPath)

	if pc.transferCliReq(ser.ReqChan, cliData, r, w) == nil {
		if pc.transferCliResp(cliData, r, w) != nil {
			//TODO: Log err
			return
		}
	} else {
		cliData.Body.Close()
		//TODO: Log err
	}
}

func (pc *ProxyCtrl) transferCliReq(reqChan chan<- *model.ProxyCliData, data *model.ProxyCliData, r *http.Request, w http.ResponseWriter) (err error) {
	select {
	case reqChan <- data:
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

func (pc *ProxyCtrl) transferCliResp(data *model.ProxyCliData, r *http.Request, w http.ResponseWriter) (err error) {
	select {
	case respData := <-data.RespChan:
		exclude := map[string]bool{
			"Httprelay-Proxy-Jobid": true,
		}
		if err = respData.Header.WriteSubset(w, exclude); err == nil {
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
