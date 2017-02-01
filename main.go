package main

import (
	"flag"
	"strconv"
	"fmt"
	"net"
	"net/http"
	"log"
	_ "net/http/pprof"

	"github.com/cloudfoundry/bosh-hm-forwarder/config"
	"github.com/cloudfoundry/bosh-hm-forwarder/forwarder"
	"github.com/cloudfoundry/bosh-hm-forwarder/handlers"
	"github.com/cloudfoundry/bosh-hm-forwarder/tcp"
	"github.com/cloudfoundry/bosh-hm-forwarder/valuemetricsender"
	"github.com/cloudfoundry/dropsonde"
	"github.com/gorilla/mux"
)

func main() {
	configFilePath := flag.String("configPath", "", "path to the configuration file")

	flag.Parse()

	conf := config.Configuration(*configFilePath)

	dropsonde.Initialize("localhost:"+strconv.Itoa(conf.MetronPort), valuemetricsender.ForwarderOrigin)

	go func() {
		err := tcp.Open(conf.IncomingPort, forwarder.StartMessageForwarder(valuemetricsender.NewValueMetricSender()))
		if err != nil {
			log.Panicln("Could not open the TCP port", err)
		}
	}()

	log.Println("Bosh HM forwarder initialized")

	infoHandler := handlers.NewInfoHandler()
	router := mux.NewRouter()
	router.Handle("/info", infoHandler).Methods("GET")

	if conf.DebugPort > 0 {
		go pprofServer(conf.DebugPort)
	}

	log.Println(fmt.Sprintf("Starting Info Server on port %d", conf.InfoPort))

	err := http.ListenAndServe(net.JoinHostPort("", fmt.Sprintf("%d", conf.InfoPort)), router)
	if err != nil {
		log.Panicln("Failed to start up alerter: ", err)
	}
}

func pprofServer(debugPort int) {
	log.Printf("Starting Pprof Server on %d\n", debugPort)
	err := http.ListenAndServe(net.JoinHostPort("localhost", fmt.Sprintf("%d", debugPort)), nil)
	if err != nil {
		log.Panicln("Pprof Server Error", err)
	}
}
