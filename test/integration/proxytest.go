package integration

import "math/rand"

type proxyTestData [][]byte

func genProxyData() (data proxyTestData) {
	data = make(proxyTestData, 100)
	for i := 0; i < len(data); i++ {
		b := make([]byte, rand.Intn(1000000))
		rand.Read(b)
		data[i] = b
	}
	return
}
