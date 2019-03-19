package server

import (
	"fmt"
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Server struct {
	net.Listener
	stopChan  chan struct{}
	outdaters []repository.Outdater
	waiters   []Waiter
}

type Args struct {
	UnixSocket string
	Addr       string
	Port       int
}

type Waiter interface {
	Wait() <-chan struct{}
}

func NewServer(args Args) (server *Server, err error) {
	server = &Server{
		stopChan: make(chan struct{}),
	}

	if args.UnixSocket == "" {
		server.Listener, err = newTcpListener(args.Addr, args.Port)
	} else {
		server.Listener, err = newUnixListener(args.UnixSocket)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, strings.NewReader("v3"))
	})

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

func (s *Server) Start() {
	go repository.Outdate(s.outdaters, time.Minute, s.stopChan)

	go func() {
		if err := http.Serve(s, nil); err != nil && s.Active() {
			log.Print("ERROR unable to serve: ", err)
		}
	}()
	log.Println("Server is listening on " + s.Addr().String())
}

func (s *Server) Stop(timeout time.Duration) {
	log.Printf("Stopping server %s...", s.Addr())
	close(s.stopChan)
	s.waitAll(timeout)
	s.Close()
	if s.Addr().Network() == "unix" {
		os.Remove(s.Addr().String())
		//syscall.Umask(0000)
	}
	log.Println("done.")
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

func newTcpListener(addr string, port int) (listener net.Listener, err error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
}

func newUnixListener(socketPath string) (listener net.Listener, err error) {
	os.Remove(socketPath)
	//syscall.Umask(0000)
	listener, err = net.Listen("unix", socketPath)
	return
}
