package model

import (
	"gitlab.com/jonas.jasas/buffreader"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	Time        time.Time
	ContentType string
	Method      string
	Query       string
	SrcIP       string
	SrcPort     string
	Content     *buffreader.BuffReader
}

func NewData(r *http.Request) *Data {
	t := time.Now()
	contentType := r.Header.Get("Content-Type")
	srcIP := r.Header.Get("X-Real-IP")
	srcPort := r.Header.Get("X-Real-Port")
	query := filterQuery(r.URL.Query())
	method := r.Method
	content := buffreader.New(r.Body)
	content.Buff()

	return &Data{t, contentType, method, query.Encode(), srcIP, srcPort, content}
}

func filterQuery(query url.Values) (filtered url.Values) {
	filtered = url.Values{}
	for k, vals := range query {
		if !strings.EqualFold(k, "wsecret") && !strings.EqualFold(k, "seqid") {
			for _, v := range vals {
				filtered.Add(k, v)
			}
		}
	}
	return
}

func (this *Data) Size() int {
	if this == nil {
		return 0
	} else {
		return len(this.ContentType) + len(this.Query) + len(this.SrcIP) + len(this.SrcPort)
	}
}

func (this *Data) Write(w http.ResponseWriter, yourTime time.Time, expose []string) (err error) {
	w.Header().Set("Content-Type", this.ContentType)
	w.Header().Set("X-Real-IP", this.SrcIP)
	w.Header().Set("X-Real-Port", this.SrcPort)
	w.Header().Set("Httprelay-Time", toUnixMilli(this.Time))
	w.Header().Set("Httprelay-Your-Time", toUnixMilli(yourTime))
	w.Header().Set("Httprelay-Method", this.Method)
	w.Header().Set("Httprelay-Query", this.Query)
	expose = append([]string{"X-Real-IP", "X-Real-Port", "Httprelay-Time", "Httprelay-Your-Time", "Httprelay-Method", "Httprelay-Query"}, expose...)
	w.Header().Set("Access-Control-Expose-Headers", strings.Join(expose, ", "))
	io.Copy(w, this.Content)
	//log.Print(n)
	return
}

func toUnixMilli(t time.Time) string {
	mills := t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	return strconv.FormatInt(mills, 10)
}
