package repository

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"sync"
)

type linkMap map[string]*model.Link

type LinkRep struct {
	linkMap
	sync.RWMutex
	stopChan <-chan struct{}
}

func NewLinkRep(stopChan <-chan struct{}) *LinkRep {
	lnk := &LinkRep{
		linkMap:  make(linkMap),
		stopChan: stopChan,
	}
	return lnk
}

func (lr *LinkRep) Read(id string, data *model.Data, closeChan <-chan struct{}) (peerData *model.Data, ok bool) {
	lr.Lock()
	link := lr.getOrCreate(id)
	lr.Unlock()

	link.AddWaiter()
	defer link.RemoveWaiter()

	select {
	case linkData := <-link.Chan():
		select {
		case linkData.BackChan <- data:
			peerData = linkData.Data
			ok = true
		case <-closeChan:
		case <-lr.stopChan:
		}
	case <-closeChan:
	case <-lr.stopChan:
	}
	return
}

func (lr *LinkRep) Write(id string, data *model.Data, wSecret string, closeChan <-chan struct{}) (getData *model.Data, ok bool, auth bool) {
	lr.Lock()
	link := lr.getOrCreate(id)
	lr.Unlock()

	if link.WAuth(wSecret) {
		auth = true
	} else {
		return
	}

	link.AddWaiter()
	defer link.RemoveWaiter()

	linkData := model.NewLinkData(data)
	select {
	case link.Chan() <- linkData:
		select {
		case getData, ok = <-linkData.BackChan:
			link.Accessed()
		case <-closeChan:
		case <-lr.stopChan:
		}
		close(linkData.BackChan)
	case <-closeChan:
	case <-lr.stopChan:
	}

	return
}

func (lr *LinkRep) getOrCreate(id string) *model.Link {
	link, exist := lr.linkMap[id]
	if !exist {
		link = model.NewLink()
		lr.linkMap[id] = link
	}

	return link
}

func (lr *LinkRep) removeOutdated() {
	lr.Lock()
	defer lr.Unlock()
	for k, v := range lr.linkMap {
		if v.Expired() {
			delete(lr.linkMap, k)
			close(v.Chan())
		}
	}
}

func (lr *LinkRep) WaiterCount() (cnt int) {
	lr.RLock()
	defer lr.RUnlock()

	for _, v := range lr.linkMap {
		cnt += v.WaiterCount()
	}
	return
}

func (lr *LinkRep) Count() int {
	lr.RLock()
	defer lr.RUnlock()
	return len(lr.linkMap)
}
