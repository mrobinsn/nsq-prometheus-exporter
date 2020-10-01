package main

import (
	"caozhipan/nsq-prometheus-exporter/controllers"
	"crypto/tls"
	"flag"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

var (
	nsqLookupdAddress = flag.String("nsq.lookupd.address", "127.0.0.1:4161", "nsqllookupd address list with comma")
	nsqdScheme        = flag.String("nsq.scheme", "http", "scheme for talking to nsqd over http (http/https)")
)

func main() {
	flag.Parse()

	controllers.NSQDScheme = *nsqdScheme

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			controllers.SyncNodeList(*nsqLookupdAddress)
			<-ticker.C
		}
	}()

	// support self-signed certs, this is just for stats so not totally insane
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	prometheus.MustRegister(controllers.Collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9527", nil))

}
