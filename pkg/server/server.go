package server

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"net"
	"net/http"
	"time"
)

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
	http.HandleFunc("/sync/", syncCtrl.Conduct)

	linkRep := repository.NewLinkRep(server.stopChan)
	linkCtrl := controller.NewLinkCtrl(linkRep, server.stopChan)
	http.HandleFunc("/link/", linkCtrl.Conduct)

	mcastRep := repository.NewMcastRep(server.stopChan)
	mcastCtrl := controller.NewMcastCtrl(mcastRep, server.stopChan)
	http.HandleFunc("/mcast/", mcastCtrl.Conduct)

	server.outdaters = []repository.Outdater{linkRep, mcastRep}
	server.waiters = []Waiter{syncCtrl, linkCtrl, mcastCtrl}

	return
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
