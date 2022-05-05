package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"

	"github.com/rkorkosz/web"
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
	srv := web.Server(
		web.WithAddr(config.Addr),
		web.WithHandler(MultiHostProxy(hostStorage, transport)),
		web.WithTLSConfig(config.tlsConfig()),
	)
	web.RunServer(context.Background(), srv)
}
