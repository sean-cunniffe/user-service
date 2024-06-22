package main

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	createInstances()
	parentCtx := context.Background()

	log.Info("starting health probe")
	startProbes()

	serverManager.OnError(func(err error) {
		probe.SetUnReady()
	})
	serverManager.OnServing(func() {
		probe.SetReady()
	})
	log.Info("starting grpc server")
	serverManager.StartGrpcServers(parentCtx)

	log.Info("waiting for server to start")
	serverManager.ServersStarted()
	<-parentCtx.Done()
}

func startProbes() {
	healthServer := http.Server{
		Addr:    ":8080",
		Handler: &probe,
	}
	err := healthServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start health probe server: %v", err)
	}
}
