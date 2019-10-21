package controller

import (
	"net"
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

func clientIp(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		ips := strings.Split(fwd, ",")
		return ips[0]
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
