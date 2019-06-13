package server

//package main
//
//import (
//	"fmt"
//	"github.com/duoke/base/server"
//	"net/http"
//	"time"
//)
//
//func main() {
//	host := "localhost"
//	port := 8080
//	stopSignal := make(chan struct{})
//
//	graceServer := server.NewGraceExitServer(host, port, stopSignal)
//	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
//		if req.URL.Path == "/stop" {
//			close(stopSignal)
//		} else {
//			<-time.After(time.Second * 10)
//			resp.Write([]byte("ok"))
//		}
//	})
//
//	err := graceServer.ListenAndServe()
//	if err != nil {
//		fmt.Println("Err:", err)
//	}
//}

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

// GraceExitError is the error stands for the grace exit of server
var GraceExitError = errors.New("server gracefully exit")

// GraceExitServer is a http server that supports grace exit on signal
type GraceExitServer struct {
	*http.Server
	stopSignal <-chan struct{} // signal to tell the server to exit
}

// NewGraceExitServer creates a new grace exit server object with DefaultServeMux handler
func NewGraceExitServer(host string, port int, stopSignal <-chan struct{}) *GraceExitServer {
	svr := &GraceExitServer{
		stopSignal: stopSignal,
		Server:     &http.Server{},
	}
	endPoint := fmt.Sprintf("%s:%d", host, port)

	// set Addr and Handler
	svr.Addr = endPoint
	svr.Handler = nil
	return svr
}

// NewGraceExitServerWithHandler creates a new grace exit server object with self defined handler
func NewGraceExitServerWithHandler(host string, port int, stopSignal <-chan struct{}, handler http.Handler) *GraceExitServer {
	svr := &GraceExitServer{
		stopSignal: stopSignal,
		Server:     &http.Server{},
	}
	endPoint := fmt.Sprintf("%s:%d", host, port)

	// set Addr and Handler
	svr.Addr = endPoint
	svr.Handler = handler
	return svr
}

// ListenAndServe serve the http endpoint
func (srv *GraceExitServer) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	shutdownDone := make(chan struct{})

	// check the stop signal to quit tcp Accept in time
	go func() {
		defer close(shutdownDone) // notify the shutdownDone channel

		stopSignalPollInterval := time.Millisecond * 100
		ticker := time.NewTicker(stopSignalPollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-srv.stopSignal:
				srv.Shutdown(context.Background())
				http.Head(fmt.Sprintf("http://%s", addr))
				return
			default:
				<-ticker.C // block for next ticker
			}
		}
	}()

	// serve http service on address
	err = srv.Serve(GraceExitListener{
		ln.(*net.TCPListener),
		srv.stopSignal,
	})

	// wait for the listen fully exits
	<-shutdownDone

	// reset err when gracefully exits
	if err == http.ErrServerClosed || err == GraceExitError {
		err = nil
	}

	return err
}

// GraceExitListener exits the http server gracefully, which reference the implementation of tcpKeepAliveListener in net/http/server.go
type GraceExitListener struct {
	*net.TCPListener
	stopSignal <-chan struct{}
}

// Accept accepts the tcp connections
func (ln GraceExitListener) Accept() (net.Conn, error) {
	select {
	case <-ln.stopSignal:
		err := GraceExitError
		return nil, err
	default:
		tc, err := ln.AcceptTCP()
		if err != nil {
			return nil, err
		}
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(3 * time.Minute)
		return tc, nil
	}
}
