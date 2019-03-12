package main

import (
	"flag"
	"gitlab.com/jonas.jasas/httprelay/pkg/server"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var args struct {
	serverArgs server.Args
}

func init() {
	flag.StringVar(&args.serverArgs.UnixSocket, "u", "", "Bind Unix socket path")
	flag.StringVar(&args.serverArgs.Addr, "a", "", "Bind address")
	flag.IntVar(&args.serverArgs.Port, "p", 8800, "Bind port")
	flag.Parse()
}

func main() {
	log.SetPrefix(strconv.Itoa(os.Getpid()) + " ")
	log.Println("========================================================================1")
	log.Println("Starting httprelay...")

	server, _ := server.NewServer(args.serverArgs)
	server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.Stop(time.Second)
}
