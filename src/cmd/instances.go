package main

import (
	"user-service/src/configuration"
	probes "user-service/src/probes"
	"user-service/src/servers"
	"user-service/src/servers/manager"
	userservice "user-service/src/services"
)

var (
	serverManager manager.ServerManager
	probe         probes.Probes
)

func createInstances() {
	config := configuration.ReadEnvConfig()

	// Create services
	userService := userservice.NewUserService()

	// Create servers
	userServiceServer := servers.CreateUserServiceServer(userService)

	serverManager = manager.CreateServerManager(config.ServerManagerConfig, userServiceServer)

}
