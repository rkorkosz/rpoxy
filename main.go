package main

import (
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
	srv := &http.Server{
		TLSConfig: config.tlsConfig(),
		Handler:   MultiHostProxy(hostStorage),
		Addr:      config.Addr,
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
