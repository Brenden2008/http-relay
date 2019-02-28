package integration

import (
	"fmt"
	"math/rand"
	"testing"
)

type syncReqData struct {
	AData []byte
	BData []byte
}

func NewSyncData(n int) map[string]*syncReqData {
	syncMap := make(map[string]*syncReqData)
	for i := 0; i < n; i++ {
		d := syncReqData{
			make([]byte, rand.Intn(100000)),
			make([]byte, rand.Intn(100000)),
		}
		rand.Read(d.AData)
		rand.Read(d.BData)
		syncMap[genId(10)] = &d
	}
	return syncMap
}

func TestSync(t *testing.T) {
	syncReqDataMap := NewSyncData(10000)

	resChan := make(chan bool)
	for k, v := range syncReqDataMap {
		go func(k string, v *syncReqData) {
			resChan <- doReqPair("POST", "POST", fmt.Sprintf("/sync/%s", k), fmt.Sprintf("/sync/%s", k), v.AData, v.BData)
		}(k, v)
	}

	for i := 0; i < len(syncReqDataMap); i++ {
		if !<-resChan {
			t.Fail()
		}
	}
}
