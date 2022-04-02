package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/rkorkosz/web"
	"gopkg.in/yaml.v3"
)

func main() {
	conf := flag.String("c", "rpoxy.yaml", "config file")
	flag.Parse()
	config, err := InitConfig(*conf)
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		TLSConfig: config.tlsConfig,
		Handler:   MultiHostProxy(config.Hosts),
		Addr:      config.Addr,
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))
}

func MultiHostProxy(targets map[string]URL) *httputil.ReverseProxy {
	director := func(r *http.Request) {
		target := targets[r.Host]
		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		if _, ok := r.Header["User-Agent"]; !ok {
			r.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
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

type Config struct {
	Addr      string
	Email     string
	Hosts     map[string]URL
	Cert, Key string

	tlsConfig *tls.Config
}

func InitConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var conf Config
	err = yaml.NewDecoder(f).Decode(&conf)
	if err != nil {
		return nil, err
	}
	if conf.Email != "" {
		var hosts []string
		for host := range conf.Hosts {
			hosts = append(hosts, host)
		}
		conf.tlsConfig = web.AutoCertTLSConfig(conf.Email, hosts...)
	} else {
		conf.tlsConfig = web.LocalTLSConfig(conf.Cert, conf.Key)
	}
	err = validateConfig(&conf)
	return &conf, err
}

func validateConfig(conf *Config) error {
	if conf.Email == "" && conf.Cert == "" {
		return errors.New("you need to provide either acme or local config")
	}
	return nil
}
