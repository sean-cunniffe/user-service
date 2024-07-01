package main

import (
	"context"
	"net/http"
	"user-service/servers/manager"

	log "github.com/sirupsen/logrus"
)

func main() {
	createInstances()
	parentCtx := context.Background()

	log.Info("starting health probe")
	go startProbes()

	setupServerManagerCallbacks()
	log.Info("starting grpc server")
	serverManager.StartManager(parentCtx)
	<-parentCtx.Done()
}

func setupServerManagerCallbacks() {
	serverManager.OnStopServing(func(so manager.ServerOptions, err error) {
		log.Errorf("server stopped serving, %+v with error %v", so, err)
		probe.SetUnReady()
	})
	serverManager.OnServing(func() {
		log.Debugf("servers are serving")
		probe.SetReady()
	})
}

func startProbes() {
	healthServer := http.Server{
		Addr:    ":8080",
		Handler: probe,
	}
	err := healthServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start health probe server: %v", err)
	}
}
