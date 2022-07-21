package web

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"log"
	"net/url"
	"time"

	"github.com/johanbrandhorst/certify"
	"github.com/johanbrandhorst/certify/issuers/vault"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// LocalAndAutoCert is loading both local and autocert Certificates
// Local certificates are checked first
func LocalAndAutoCert(cert, key, email string, policy autocert.HostPolicy) *tls.Config {
	localConf := LocalTLSConfig(cert, key)
	autoConf := TlsWithAutoCert(localConf, email, policy)
	conf := autoConf.Clone()
	conf.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		for _, cert := range conf.Certificates {
			if err := clientHello.SupportsCertificate(&cert); err == nil {
				return &cert, nil
			}
		}
		return autoConf.GetCertificate(clientHello)
	}
	return conf
}

func AutoCertTLSConfig(email string, policy autocert.HostPolicy) *tls.Config {
	return TlsWithAutoCert(BaseTLSConfig(), email, policy)
}

func AutoCertWhitelist(email string, hosts ...string) *tls.Config {
	return AutoCertTLSConfig(email, autocert.HostWhitelist(hosts...))
}

func TlsWithAutoCert(conf *tls.Config, email string, policy autocert.HostPolicy) *tls.Config {
	conf = conf.Clone()
	m := &autocert.Manager{
		Cache:      autocert.DirCache("./letsencrypt/"),
		Prompt:     autocert.AcceptTOS,
		Email:      email,
		HostPolicy: policy,
	}
	conf.NextProtos = []string{
		"h3", "h2", "http/1.1", // enable HTTP/2
		acme.ALPNProto, // enable tls-alpn ACME challenges
	}
	conf.GetCertificate = m.GetCertificate
	return conf
}

func TlsWithLocalCert(conf *tls.Config, certFile, keyFile string) (*tls.Config, error) {
	conf = conf.Clone()
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	if conf.Certificates == nil {
		conf.Certificates = []tls.Certificate{}
	}
	conf.Certificates = append(conf.Certificates, cert)
	return conf, nil
}

func LocalTLSConfig(certFile, keyFile string) *tls.Config {
	conf, err := TlsWithLocalCert(BaseTLSConfig(), certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

func VaultTLSConfig(vaultAddr, vaultToken, vaultRole, cn string, ttl time.Duration) (*tls.Config, error) {
	vaultURL, err := url.Parse(vaultAddr)
	if err != nil {
		return nil, err
	}
	issuer := &vault.Issuer{
		URL:        vaultURL,
		Token:      vaultToken,
		Role:       vaultRole,
		TimeToLive: ttl,
	}
	c := &certify.Certify{
		CommonName: cn,
		Issuer:     issuer,
		Cache:      certify.NewMemCache(),
		CertConfig: &certify.CertConfig{
			KeyGenerator: keyGeneratorFunc(func() (crypto.PrivateKey, error) {
				return rsa.GenerateKey(rand.Reader, 2048)
			}),
		},
	}
	tlsConfig := BaseTLSConfig()
	tlsConfig.GetCertificate = c.GetCertificate
	return tlsConfig, nil
}

type keyGeneratorFunc func() (crypto.PrivateKey, error)

func (kgf keyGeneratorFunc) Generate() (crypto.PrivateKey, error) {
	return kgf()
}

func BaseTLSConfig() *tls.Config {
	return &tls.Config{
		NextProtos: []string{
			"h2", "http/1.1",
		},
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
