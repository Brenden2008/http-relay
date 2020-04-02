package model

type dataMap map[string]*ProxyData

// Creating when first request to/from server is made
type ProxySer struct {
	reqChan chan *ProxyData
	jobMap  dataMap
	comm
}

func NewProxySer() *ProxySer {
	return &ProxySer{
		jobMap: make(dataMap),
		comm:   newComm(),
	}
}

func (ps *ProxySer) AddReq(proxyData *ProxyData, closeChan <-chan struct{}) {
	select {
	case ps.reqChan <- proxyData:
	case <-closeChan:
	}
}

func (ps *ProxySer) TakeReq(closeChan <-chan struct{}) (proxyData *ProxyData) {
	select {
	case proxyData = <-ps.reqChan:
	case <-closeChan:
	}
	return
}

func (ps *ProxySer) AddJob(jobId string, proxyData *ProxyData) {
	ps.Lock()
	defer ps.Unlock()
	ps.jobMap[jobId] = proxyData
}

func (ps *ProxySer) TakeJob(jobId string) (proxyData *ProxyData, ok bool) {
	ps.Lock()
	defer ps.Unlock()
	if proxyData, ok = ps.jobMap[jobId]; ok {
		delete(ps.jobMap, jobId)
	}
	return
}
