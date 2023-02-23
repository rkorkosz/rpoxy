package main

import (
	"context"
	"crypto/tls"
	"errors"
	"os"

	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/yaml.v3"
)

type Config struct {
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
	err = validateConfig(&conf)
	return &conf, err
}

func (c *Config) GetHost(_ context.Context, host string) (URL, error) {
	return c.Hosts[host], nil
}

func (c *Config) tlsConfig() *tls.Config {
	hosts := []string{}
	for h := range c.Hosts {
		hosts = append(hosts, h)
	}
	if c.Email != "" && c.Cert != "" {
		return LocalAndAutoCert(c.Cert, c.Key, c.Email, autocert.HostWhitelist(hosts...))
	}
	if c.Cert != "" {
		return LocalTLSConfig(c.Cert, c.Key)
	}
	if c.Email != "" {
		return AutoCertWhitelist(c.Email, hosts...)
	}
	return nil
}

func validateConfig(conf *Config) error {
	if conf.Email == "" && conf.Cert == "" {
		return errors.New("you need to provide either acme or local config")
	}
	return nil
}
