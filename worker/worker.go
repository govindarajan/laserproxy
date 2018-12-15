package worker

import (
	"crypto/tls"
	"io"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/govindarajan/laserproxy/helper"
	"github.com/govindarajan/laserproxy/logger"
)

var transports = make(map[string]http.RoundTripper)
var mutex = &sync.RWMutex{}

func getTransport(ip string) http.RoundTripper {
	// IF found in map, return.
	var trans http.RoundTripper
	var ok bool
	if trans, ok = transports[ip]; ok {
		return trans
	}

	// otherwise form, store and return
	// TODO: Get values from config
	ipaddr, err := net.ResolveTCPAddr("tcp", ip+":0")
	if err != nil {
		// Incase of error, return Default Transport.
		logger.LogWarn(err.Error())
		return http.DefaultTransport
	}
	trans = &http.Transport{
		Proxy:        http.ProxyFromEnvironment,
		MaxIdleConns: 1,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
			LocalAddr: ipaddr,
		}).DialContext,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   2 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	mutex.Lock()
	transports[ip] = trans
	mutex.Unlock()
	return trans
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {

	outgoing := getOutgoingRoute()

	resp, err := getTransport(outgoing).RoundTrip(r)
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

func StartProxy() {
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

	if e := server.ListenAndServe(); e != nil {
		logger.LogError(e.Error())
	}
}

func getOutgoingRoute() string {
	// TODO: Based on the config, return best route or wighet based route.
	ips, _ := helper.GetLocalIPs()
	r := ips[rand.Intn(len(ips))]
	logger.LogDebug("Outbound Route:" + r.IP)
	return r.IP
}

func getTargetIPIfAny() *string {
	return nil
}

func Start() {
	// Get Ougoing Route
	// Get Target IP if there is any
	// Make request with above details.
}
