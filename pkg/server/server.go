package server

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"net"
	"net/http"
	"strings"
	"time"
)

var Version string = "test"

type Server struct {
	net.Listener
	stopChan  chan struct{}
	errChan   chan error
	outdaters []repository.Outdater
	waiters   []Waiter
}

type Waiter interface {
	Wait() <-chan struct{}
}

func NewServer(listener net.Listener) (server *Server) {
	server = &Server{
		stopChan: make(chan struct{}),
		errChan:  make(chan error, 1),
	}

	server.Listener = listener

	syncRep := repository.NewSyncRep(server.stopChan)
	syncCtrl := controller.NewSyncCtrl(syncRep, server.stopChan)
	http.HandleFunc("/sync/", corsHandler(syncCtrl.Conduct, []string{}))

	linkRep := repository.NewLinkRep(server.stopChan)
	linkCtrl := controller.NewLinkCtrl(linkRep, server.stopChan)
	http.HandleFunc("/link/", corsHandler(linkCtrl.Conduct, []string{}))

	mcastRep := repository.NewMcastRep(server.stopChan)
	mcastCtrl := controller.NewMcastCtrl(mcastRep, server.stopChan)
	http.HandleFunc("/mcast/", corsHandler(mcastCtrl.Conduct, []string{"Httprelay-Seqid"}))

	server.outdaters = []repository.Outdater{linkRep, mcastRep}
	server.waiters = []Waiter{syncCtrl, linkCtrl, mcastCtrl}

	return
}

func corsHandler(h http.HandlerFunc, expose []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		cors(w, r, expose)
		if r.Method != "OPTIONS" {
			h(w, r)
		}
	}
}

func cors(w http.ResponseWriter, r *http.Request, expose []string) {
	w.Header().Set("Httprelay-Version", Version)

	if r.Method == "OPTIONS" {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	} else {
		expose = append(expose, "Content-Length, X-Real-IP, X-Real-Port, Httprelay-Version, Httprelay-Time, Httprelay-Your-Time, Httprelay-Method, Httprelay-Query")
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(expose, ", "))
	}
}

func (s *Server) Start() <-chan error {
	go repository.Outdate(s.outdaters, time.Minute, s.stopChan)

	go func() {
		if err := http.Serve(s, nil); err != nil && s.Active() {
			s.Stop(time.Second)
			s.errChan <- err
		}
	}()
	return s.errChan
}

func (s *Server) Stop(timeout time.Duration) {
	close(s.stopChan)
	s.waitAll(timeout)
	s.Close()
}

func (s *Server) Active() bool {
	select {
	case <-s.stopChan:
		return false
	default:
		return true
	}
}

func (s *Server) waitAll(timeout time.Duration) {
	t := time.NewTimer(timeout)
	for _, w := range s.waiters {
		select {
		case <-w.Wait():
		case <-t.C:
		}
	}
	t.Stop()
}
