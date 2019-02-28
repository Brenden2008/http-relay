package controller

import (
	"net/http"
	"strings"
)

func wSecret(r *http.Request) string {
	for k, v := range r.URL.Query() {
		if strings.EqualFold(k, "wsecret") {
			return v[0]
		}
	}

	return r.Header.Get("Httprelay-WSecret")
}
