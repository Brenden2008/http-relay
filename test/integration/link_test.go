package integration

//import (
//	"fmt"
//	"math/rand"
//	"testing"
//)
//
//type linkReqData struct {
//	AData []byte
//	BData []byte
//}
//
//func NewSyncData(n int) map[string]*[] {
//	syncMap := make(map[string]*syncReqData)
//	for i := 0; i < n; i++ {
//		d := syncReqData{
//			make([]byte, rand.Intn(100000)),
//			make([]byte, rand.Intn(100000)),
//		}
//		rand.Read(d.AData)
//		rand.Read(d.BData)
//		syncMap[genId(10)] = &d
//	}
//	return syncMap
//}
//
//
//func TestLink(t *testing.T) {
//
//	genId(10)
//
//	resChan := make(chan bool)
//	for k, v := range syncReqDataMap {
//		go func(k string, v *syncReqData) {
//			resChan <- doReqPair("POST", "GET", fmt.Sprintf("/sync/%s", k), fmt.Sprintf("/sync/%s", k), v.AData, v.BData)
//		}(k, v)
//	}
//
//	for i := 0; i < len(syncReqDataMap); i++ {
//		if !<-resChan {
//			t.Fail()
//		}
//	}
//}
