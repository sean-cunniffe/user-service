package main

import (
	"user-service/src/configuration"
	probes "user-service/src/probes"
	"user-service/src/servers"
	"user-service/src/servers/manager"
	userservice "user-service/src/services"

	log "github.com/sirupsen/logrus"
)

var (
	serverManager manager.ServerManager
	probe         *probes.Probes
)

func createInstances() {
	config := configuration.ReadEnvConfig()
	log.Infof("creating instances with config %+v", config)

	// Create services
	userService := userservice.NewUserService()

	// Create servers
	userServiceServer := servers.CreateUserServiceServer(userService)

	serverManager = manager.CreateServerManager(config.ServerManagerConfig, userServiceServer)

	// create probes
	probe = probes.NewProbes(config.ProbeConfig)
}
