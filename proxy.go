package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type HostStorage interface {
	GetHost(context.Context, string) (URL, error)
}

func MultiHostProxy(targets HostStorage) *httputil.ReverseProxy {
	director := func(r *http.Request) {
		target, err := targets.GetHost(r.Context(), r.Host)
		if err != nil {
			log.Println(err)
			return
		}
		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		if _, ok := r.Header["User-Agent"]; !ok {
			r.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director, Transport: httpTransport()}
}

type URL url.URL

func (u *URL) UnmarshalText(b []byte) error {
	parsed, err := url.Parse(string(b))
	if err != nil {
		return err
	}
	*u = URL(*parsed)
	return nil
}

func httpTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
