package repository

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"sync"
)

type serMap map[string]*model.ProxySer

type ProxyRep struct {
	serMap
	sync.Mutex
}

func NewProxy() *ProxyRep {
	return &ProxyRep{
		serMap: make(serMap),
	}
}

func (pr *ProxyRep) GetServer(serId string) *model.ProxySer {
	pr.Lock()
	defer pr.Unlock()
	proxySer, ok := pr.serMap[serId]
	if !ok {
		proxySer = model.NewProxySer()
		pr.serMap[serId] = proxySer
	}

	return proxySer
}
