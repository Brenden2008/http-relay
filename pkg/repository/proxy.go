package repository

import "gitlab.com/jonas.jasas/httprelay/pkg/model"

type proxyReqMap map[string]*model.Proxy
type proxyMap map[string]*proxyReqMap

type ProxyRep struct {
}
