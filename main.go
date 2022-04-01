package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rkorkosz/web"
)

func main() {
	addr := flag.String("bind", ":8001", "bind address")
	cert := flag.String("cert", "cert.pem", "certificate")
	key := flag.String("key", "key.pem", "certificate key")
	email := flag.String("email", "", "acme email")
	host := flag.String("host", "", "host for certificate")
	flag.Parse()
	target, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	var tlsConfig *tls.Config
	if *cert != "" {
		tlsConfig = web.LocalTLSConfig(*cert, *key)
	}
	if *email != "" {
		tlsConfig = web.AutoCertTLSConfig(*email, *host)
	}
	srv := &http.Server{
		TLSConfig: tlsConfig,
		Handler:   httputil.NewSingleHostReverseProxy(target),
		Addr:      *addr,
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
