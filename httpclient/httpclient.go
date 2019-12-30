package httpclient

import (
	"net"
	"net/http"
	"time"
)

var (
	defaultTransport http.RoundTripper = &http.Transport{
		// Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			// dial timeout
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          500,
		IdleConnTimeout:       360 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,

		// per host
		MaxConnsPerHost:     200,
		MaxIdleConnsPerHost: 50,
	}

	client = &http.Client{
		Transport: defaultTransport,
		// connect + req + resp
		Timeout: 10 * time.Second,
	}
)

func init() {
	http.DefaultClient = client
}

func SetTimeout(t time.Duration) {
	client = &http.Client{
		Transport: defaultTransport,
		// connect + req + resp
		Timeout: 10 * time.Second,
	}
	http.DefaultClient = client
}
