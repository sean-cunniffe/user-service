package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func init() {
	createInstances()
}

func main() {
	parentCtx := context.Background()
	log.Info("starting grpc server")
	serverManager.StartGrpcServers(parentCtx)
	log.Info("waiting for server to start")
	serverManager.ServersStarted()
	<-parentCtx.Done()
}
