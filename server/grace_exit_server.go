package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// GraceExitServer is a http server supports grace exit
type GraceExitServer struct {
	*http.Server
	stopSignal   <-chan struct{} // signal to tell the server to exit
	notifySignal chan struct{}   // signal to notify when server gracelly exited
}

// NewGraceExitServer creates a new server object
func NewGraceExitServer(host string, port int, stopSignal <-chan struct{}, notifySignal chan struct{}) *GraceExitServer {
	svr := &GraceExitServer{
		stopSignal:   stopSignal,
		notifySignal: notifySignal,
		Server:       &http.Server{},
	}
	endPoint := fmt.Sprintf("%s:%d", host, port)
	svr.Addr = endPoint
	svr.Handler = nil
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

	// check the stop signal to quit accept in time
	go func() {
		for {
			select {
			case <-srv.stopSignal:
				http.Get(fmt.Sprintf("http://%s", addr))
			default:
				<-time.After(time.Millisecond * 100)
			}
		}
	}()

	// serve http service on address
	err = srv.Serve(GraceExitListener{
		ln.(*net.TCPListener),
		srv.stopSignal,
	})

	// shutdown the server after accept quits
	err = srv.Shutdown(context.Background())
	if err != nil {
		return err
	}

	// fire notify signal
	close(srv.notifySignal)

	// if grace exit, reset error
	if _, ok := err.(GraceExitError); ok {
		return nil
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
		err := GraceExitError{}
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

// GraceExitError is the error stands for the grace exit of server
type GraceExitError struct {
}

func (e GraceExitError) Error() string {
	return "server gracefully exit"
}
