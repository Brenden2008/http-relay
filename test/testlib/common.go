package testlib

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func newReq(method string, url string, header map[string]string, data string) (r *http.Request) {
	r, _ = http.NewRequest(method, url, strings.NewReader(data))
	if header != nil {
		for k, v := range header {
			r.Header.Add(k, v)
		}
	}
	return
}

func RespDataEq(body io.Reader, data string) bool {
	respData, _ := ioutil.ReadAll(body)
	return string(respData) == data
}
