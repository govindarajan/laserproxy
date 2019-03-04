package worker

import (
	"crypto/tls"
	"database/sql"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/monitor"
	"github.com/govindarajan/laserproxy/store"
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

func handleHTTP(w http.ResponseWriter, r *http.Request, id int) {

	outgoing := getOutgoingRoute()
	target := getTargetIPIfAny(r.URL.Host)
	if target != nil {
		r.URL.Host = *target
	}
	resp, err := getTransport(outgoing).RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	resp.Header.Add("X-Proxy", "LaserProxy")
	// TODO: Should we retry incase of timeout??
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

	//call monitor to do healthchecks
	//Todo: fix go routine leaks here
	// monitor, err := monitor.NewMonitor()
	// if err != nil {
	// 	logger.LogCritical("unable to do healthchecks")
	// }
	// monitor.Schedule(10000)
	server := &http.Server{
		Addr: ":8888",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				// TLS:

			} else {
				handleHTTP(w, r, 0)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	if e := server.ListenAndServe(); e != nil {
		logger.LogError(e.Error())
	}
	logger.LogInfo("Starting Frontends")
	StartFrontEnds(maindb)
}

func getOutgoingRoute() string {
	lrs, err := store.ReadLocalRoutes(maindb)
	if err != nil {
		logger.LogError("GetOBRoute: " + err.Error())
		return ""
	}
	route := lrs[rand.Intn(len(lrs))]
	return route.IP.String()
}

func getTargetIPIfAny(host string) *string {

	return nil
}

// StartFrontEnds used to start all the front end proxies
// by reading the frondends table.
func StartFrontEnds(db *sql.DB) error {
	fes, err := store.ReadFrontends(db)
	if err != nil {
		return err
	}
	for _, fe := range fes {
		// Start proxies
		logger.LogDebug(fe.ListenAddr.String())
		server, err := startProxy(&fe)
		if err != nil {
			continue
		}
		frontends[fe.Id] = server
	}
	return nil
}

func startProxy(fe *store.Frontend) (*http.Server, error) {
	server := &http.Server{
		Addr: fe.ListenAddr.String() + ":" + strconv.Itoa(fe.Port),

		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	if fe.Type == store.PrTypeForward {
		// Forward proxy
		server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				// TLS:

			} else {
				handleHTTP(w, r, fe.Id)
			}
		})
	} else {
		// Reverse proxy
		server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleReverseProxyReq(w, r, fe)
		})
	}

	err := server.ListenAndServe()
	return server, err
}

func handleReverseProxyReq(w http.ResponseWriter, r *http.Request, fe *store.Frontend) {
	bends := monitor.GetHealthyBackends(maindb, fe)
	for _, be := range bends {
		purl, err := url.Parse("http://" + be.Host)
		if err != nil {
			logger.LogError("ReverseProxy: Readbackeds " + err.Error())
			http.Error(w, "", http.StatusServiceUnavailable)
			continue
		}
		// create reverse proxy.
		// TODO: Optimize here
		proxy := httputil.NewSingleHostReverseProxy(purl)
		// Set the transport to
		proxy.Transport = getTransport(getOutgoingRoute())

		// Update the headers to allow for SSL redirection
		r.URL.Host = purl.Host
		r.URL.Scheme = purl.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = purl.Host

		proxy.ServeHTTP(w, r)
		// TODO: Should we retry on 5xx???
		break
	}
}
