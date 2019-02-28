package integration

import (
	"bytes"
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/rwmock"
	"golang.org/x/net/context"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

func genId(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func doReq(method, path string, body io.Reader) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, body)

	cancelChan := make(chan struct{})
	srvState := controller.NewSrvState()
	reqCtx := controller.NewReqCtx(cancelChan, srvState)
	ctx := req.Context()
	ctx = context.WithValue(ctx, controller.ReqCtxKey, reqCtx)
	req = req.WithContext(ctx)

	//controller.Sync(resp, req)
	handler := http.HandlerFunc(controller.Sync)
	handler.ServeHTTP(resp, req)

	return resp
}

func doReqPair(methodA, methodB, pathA, pathB string, bodyA, bodyB []byte) bool {

	respAChan := make(chan *httptest.ResponseRecorder)

	rA := rwmock.NewShaperRand(bytes.NewReader(bodyA), 1, 1000, 0, time.Microsecond)
	rB := rwmock.NewShaperRand(bytes.NewReader(bodyB), 1, 1000, 0, time.Microsecond)

	go func() {
		time.Sleep(time.Duration(rand.Int63n(1000000)))
		respAChan <- doReq(methodA, pathA, rA)
	}()

	time.Sleep(time.Duration(rand.Int63n(1000000)))
	respB := doReq(methodB, pathB, rB)

	respA := <-respAChan

	a := respA.Body.Bytes()
	b := respB.Body.Bytes()

	return bytes.Compare(a, bodyB) == 0 && bytes.Compare(b, bodyA) == 0
}
