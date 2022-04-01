package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/lucas-clemente/quic-go/qlog"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func Server(opts ...func(*http.Server)) *http.Server {
	srv := &http.Server{
		Addr:         ":443",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

func WithAddr(addr string) func(*http.Server) {
	return func(srv *http.Server) {
		srv.Addr = addr
	}
}

func WithTLSConfig(conf *tls.Config) func(*http.Server) {
	return func(srv *http.Server) {
		srv.TLSConfig = conf
	}
}

func WithHandler(h http.Handler) func(*http.Server) {
	return func(srv *http.Server) {
		srv.Handler = h
	}
}

func RunServer(ctx context.Context, srv *http.Server) {
	go func() {
		if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	go func() {
		s := http3.Server{
			Server: srv,
			QuicConfig: &quic.Config{
				Tracer: qlog.NewTracer(func(_ logging.Perspective, connectionID []byte) io.WriteCloser {
					return os.Stdout
				}),
			},
		}

		if err := s.ListenAndServeTLS("", ""); err != nil {
			log.Fatal(err)
		}
	}()
	log.Print("Started")
	<-ctx.Done()
	log.Print("Stopping")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
	log.Print("Stopped")
}

func HTTP3Middleware(next http.Handler, addr string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("alt-svc", fmt.Sprintf(`h3="%s"; ma=2592000,h3-29="%s"; ma=2592000,quic="%s"; ma=2592000; v="46,43"`, addr, addr, addr))
		next.ServeHTTP(w, r)
	})
}

func AutoCertTLSConfig(email, host string) *tls.Config {
	m := &autocert.Manager{
		Cache:      autocert.DirCache("./letsencrypt/"),
		Prompt:     autocert.AcceptTOS,
		Email:      email,
		HostPolicy: autocert.HostWhitelist(host),
	}
	conf := baseTLSConfig()
	conf.NextProtos = []string{
		"h3", "h2", "http/1.1", // enable HTTP/2
		acme.ALPNProto, // enable tls-alpn ACME challenges
	}
	conf.GetCertificate = m.GetCertificate
	srv := Server(WithAddr(":http"))
	srv.Handler = m.HTTPHandler(nil)
	go srv.ListenAndServe()
	return conf
}

func LocalTLSConfig(certFile, keyFile string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	conf := baseTLSConfig()
	conf.Certificates = []tls.Certificate{cert}
	return conf
}

func baseTLSConfig() *tls.Config {
	return &tls.Config{
		NextProtos: []string{
			"h2", "http/1.1",
		},
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
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
