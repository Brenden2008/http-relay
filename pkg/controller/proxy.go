package controller

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
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