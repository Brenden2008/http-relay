package model

import "sync"

type dataMap map[string]*ProxyData

// Creating when first request to/from server is made
type ProxySer struct {
	ReqChan chan *ProxyData
	jobMap  dataMap
	jobMapL sync.Mutex

	comm comm
}

func NewProxySer() *ProxySer {
	return &ProxySer{
		ReqChan: make(chan *ProxyData),
		jobMap:  dataMap{},
		comm:    newComm(),
	}
}

func (ps *ProxySer) AddJob(jobId string, proxyData *ProxyData) {
	ps.jobMapL.Lock()
	defer ps.jobMapL.Unlock()
	ps.jobMap[jobId] = proxyData
}

func (ps *ProxySer) TakeJob(jobId string) (proxyData *ProxyData, ok bool) {
	ps.jobMapL.Lock()
	defer ps.jobMapL.Unlock()
	if proxyData, ok = ps.jobMap[jobId]; ok {
		delete(ps.jobMap, jobId)
	}
	return
}
