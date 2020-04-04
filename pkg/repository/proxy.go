package repository

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"sync"
)

type serMap map[string]*model.ProxySer

type ProxyRep struct {
	serMap  serMap
	serMapL sync.Mutex
}

func NewProxyRep() *ProxyRep {
	return &ProxyRep{
		serMap: serMap{},
	}
}

func (pr *ProxyRep) GetSer(serId string) *model.ProxySer {
	pr.serMapL.Lock()
	defer pr.serMapL.Unlock()
	proxySer, ok := pr.serMap[serId]
	if !ok {
		proxySer = model.NewProxySer()
		pr.serMap[serId] = proxySer
	}
	return proxySer
}
