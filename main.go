package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
)

func main() {
	conf := flag.String("c", "rpoxy.yaml", "config file")
	flag.Parse()
	config, err := InitConfig(*conf)
	if err != nil {
		log.Fatal(err)
	}
	var hostStorage HostStorage
	hostStorage = config
	if config.KV != "" {
		hostStorage, err = NewKV(config.KV, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	transport := httpTransport()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	srv := http.Server{
		Addr:      config.Addr,
		Handler:   MultiHostProxy(hostStorage, transport),
		TLSConfig: config.tlsConfig(),
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
