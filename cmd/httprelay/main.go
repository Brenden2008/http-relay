package main

import (
	"../../pkg/server"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var args struct {
	Addr   string
	Port   int
	Socket string
}

func init() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	flag.StringVar(&args.Addr, "a", "", "Bind address")
	flag.IntVar(&args.Port, "p", port, "Bind port")
	flag.StringVar(&args.Socket, "u", "", "Bind Unix socket path")
	flag.Parse()
}

func listener() (net.Listener, error) {
	if args.Socket == "" {
		return net.Listen("tcp", fmt.Sprintf("%s:%d", args.Addr, args.Port))
	} else {
		os.Remove(args.Socket)
		//syscall.Umask(0000)
		return net.Listen("unix", args.Socket)
	}
}

func main() {
	fmt.Println("========================================================================")
	fmt.Println("Starting httprelay...")

	if listener, err := listener(); err == nil {
		server := server.NewServer(listener)
		errChan := server.Start()

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("health req")
			io.Copy(w, strings.NewReader("v22"))
		})

		fmt.Println("Server is listening on " + server.Addr().String())

		intChan := make(chan os.Signal, 1)
		signal.Notify(intChan, os.Interrupt)

		select {
		case <-intChan:
			fmt.Printf("Stopping server %s...", server.Addr())
			server.Stop(time.Second)
		case err := <-errChan:
			fmt.Fprintln(os.Stderr, "ERROR unable to serve: ", err)
		}

		if server.Addr().Network() == "unix" {
			os.Remove(server.Addr().String())
			//syscall.Umask(0000)
		}
		fmt.Println("done.")
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
}
