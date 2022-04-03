package main

import (
	"context"
	"crypto/tls"
	"errors"
	"os"

	"github.com/rkorkosz/web"
	"gopkg.in/yaml.v3"
)

type Config struct {
	tlsConfig *tls.Config

	Addr      string
	Email     string
	Hosts     map[string]URL
	KV        string
	Cert, Key string
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

func (c *Config) GetHost(_ context.Context, host string) (URL, error) {
	return c.Hosts[host], nil
}

func validateConfig(conf *Config) error {
	if conf.Email == "" && conf.Cert == "" {
		return errors.New("you need to provide either acme or local config")
	}
	return nil
}
