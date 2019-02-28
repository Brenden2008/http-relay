package controller

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"net/http"
	"strings"
	"time"
)

type LinkRep interface {
	Read(id string, data *model.Data, cancelChan <-chan struct{}) (peerData *model.Data, ok bool)
	Write(id string, data *model.Data, wSecret string, cancelChan <-chan struct{}) (getData *model.Data, ok bool, auth bool)
}

type LinkCtrl struct {
	rep      LinkRep
	stopChan <-chan struct{}
	*model.Waiters
}

func NewLinkCtrl(rep LinkRep, stopChan <-chan struct{}) *LinkCtrl {
	return &LinkCtrl{
		rep:      rep,
		stopChan: stopChan,
		Waiters:  model.NewWaiters(),
	}
}

func (lc *LinkCtrl) Conduct(w http.ResponseWriter, r *http.Request) {
	lc.AddWaiter()
	defer lc.RemoveWaiter()

	defer r.Body.Close()

	select {
	case <-lc.stopChan:
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	default:
	}

	pathArr := strings.Split(r.URL.Path, "/")
	id := pathArr[len(pathArr)-1]

	data := model.NewData(r)
	yourTime := time.Now()

	if strings.EqualFold(r.Method, http.MethodGet) {
		if peerData, ok := lc.rep.Read(id, data, r.Context().Done()); ok {
			lc.AddWaiter()
			<-data.Content.Buff()
			peerData.Write(w, yourTime, nil)
			lc.RemoveWaiter()
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	} else if strings.EqualFold(r.Method, http.MethodPost) {
		if peerData, ok, auth := lc.rep.Write(id, data, wSecret(r), r.Context().Done()); ok && auth {
			<-data.Content.Buff()
			peerData.Write(w, yourTime, nil)
		} else {
			if auth {
				w.WriteHeader(http.StatusServiceUnavailable)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
