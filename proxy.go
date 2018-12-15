package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"./helper"
	"github.com/govindarajan/laserproxy/logger"
)

var addr = &net.TCPAddr{net.IP{192, 168, 168, 149}, 0, ""}
var CustomTransport = &http.Transport{
	Proxy:        http.ProxyFromEnvironment,
	MaxIdleConns: 1,
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
		LocalAddr: addr,
	}).DialContext,
	IdleConnTimeout:       10 * time.Second,
	TLSHandshakeTimeout:   2 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {

	resp, err := CustomTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	resp.Header.Add("X-Proxy", "LaserProxy")
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func main() {
	logger.LogInfo("Hello Laser")
	helper.ConfigureRoute()

	server := &http.Server{
		Addr: ":8888",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				// TLS:

			} else {
				handleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Fatal(server.ListenAndServe())
}
