package main

import (
	"user-service/configuration"
	probes "user-service/probes"
	"user-service/servers"
	"user-service/servers/manager"
	userservice "user-service/services"
)

var (
	serverManager manager.ServerManager
	probe         *probes.Probes
)

func createInstances() {
	config := configuration.ReadEnvConfig()

	// Create services
	userService := userservice.NewUserService()

	// Create servers
	userServiceServer := servers.CreateUserServiceServer(userService)

	serverManager = manager.CreateServerManager(config.ServerManagerConfig, userServiceServer)

	// create probes
	probe = probes.NewProbes(config.ProbeConfig)
}
