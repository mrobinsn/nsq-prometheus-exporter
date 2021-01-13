package main

import (
	"caozhipan/nsq-prometheus-exporter/controllers"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

var (
	nsqLookupdAddress = flag.String("nsq.lookupd.address", "127.0.0.1:4161", "nsqllookupd address list with comma")
	nsqdScheme        = flag.String("nsq.scheme", "http", "scheme for talking to nsqd over http (http/https)")
	enableTLS         = flag.Bool("tls", false, "enable mTLS with nsq")
	caRootCertPath    = flag.String("ca-root-cert-path", "", "path to the root CA certificate")
	certPath          = flag.String("cert-path", "", "path to the client cert")
	keyPath           = flag.String("key-path", "", "path to the client cert key")
)

func main() {
	flag.Parse()

	controllers.NSQDScheme = *nsqdScheme

	go func() {
		for {
			controllers.SyncNodeList(*nsqLookupdAddress)
			<-time.After(10 * time.Second)
		}
	}()

	// configure mTLS if enabled
	if *enableTLS {
		cfg := tls.Config{}
		cfg.MinVersion = tls.VersionTLS12

		// Load in the root CA cert
		rootCAs := x509.NewCertPool()
		certs, err := ioutil.ReadFile(*caRootCertPath)
		if err != nil {
			panic(errors.Errorf("failed to load root ca cert: %v", err))
		}

		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			panic(errors.New("failed to parse root ca cert"))
		}
		cfg.RootCAs = rootCAs

		// Load in the client certificate
		cert, err := tls.LoadX509KeyPair(*certPath, *keyPath)
		if err != nil {
			panic(errors.Errorf("failed to load client certificate: %v", err))
		}
		cfg.Certificates = append(cfg.Certificates, cert)

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &cfg
	}

	prometheus.MustRegister(controllers.Collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9527", nil))
}
