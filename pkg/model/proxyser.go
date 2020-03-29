package model

// Creating when first request to/from server is made
type ProxySer struct {
	comMap map[string]*ProxyData
	comm
}

func NewProxyCom() *ProxySer {

}

// Server or Client can get communication channel
func (pr *ProxySer) GetCom() *ProxySer {

}
