package manager

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"user-service/src/dto"
	pb "user-service/src/generated/userservice"
	util "user-service/src/util"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	ServerManager interface {
		StartGrpcServers(ctx context.Context)
		ServersStarted() <-chan struct{}
	}

	ServerManagerConfig struct {
		ServeRetryDelay time.Duration   `yaml:"serveRetryDelay" json:"serveRetryDelay"`
		ServerOptions   []ServerOptions `yaml:"serverOptions" json:"serverOptions"`
	}
	ServerOptions struct {
		Port              int            `yaml:"port" json:"port"`
		ConnectionTimeout time.Duration  `yaml:"connectionTimeout" json:"connectionTimeout"`
		TlsConfig         *dto.TLSConfig `yaml:"tlsConfig" json:"tlsConfig"`
		MaxRecvMsgSize    int            `yaml:"maxRecvMsgSize" json:"maxRecvMsgSize"`
	}
	serverManager struct {
		config            ServerManagerConfig
		userServiceServer pb.UserServiceServer
		wg                *sync.WaitGroup
	}
)

// CreateServerManager creates the manager that manages any servers
func CreateServerManager(config ServerManagerConfig, userServiceServer pb.UserServiceServer) ServerManager {
	util.NotNil(userServiceServer, config)
	wg := &sync.WaitGroup{}
	wg.Add(len(config.ServerOptions))
	return &serverManager{
		userServiceServer: userServiceServer,
		config:            config,
		wg:                wg,
	}
}

func getGrpcServerOptions(serverOptions ServerOptions) []grpc.ServerOption {
	options := []grpc.ServerOption{
		grpc.ConnectionTimeout(serverOptions.ConnectionTimeout),
		grpc.MaxRecvMsgSize(serverOptions.MaxRecvMsgSize),
	}
	if serverOptions.TlsConfig == nil {
		log.Printf("running server without TLS")
		return options
	}
	// Load TLS config
	tlsConfig := util.GetTLSConfig(*serverOptions.TlsConfig)
	options = append(options, grpc.Creds(credentials.NewTLS(tlsConfig)))
	return options
}

func (sm *serverManager) ServersStarted() <-chan struct{} {
	ready := make(chan struct{})
	go func() {
		sm.wg.Wait()
		close(ready)
	}()
	return ready
}

// start the grpc server from the serveroptions
func (sm *serverManager) StartGrpcServers(ctx context.Context) {
	servers := make([]*grpc.Server, len(sm.config.ServerOptions))
	for i, options := range sm.config.ServerOptions {
		grpcServerOptions := getGrpcServerOptions(options)
		server := grpc.NewServer(grpcServerOptions...)
		pb.RegisterUserServiceServer(server, sm.userServiceServer)
		servers[i] = server
		context.AfterFunc(ctx, func() {
			log.Debugf("shutting down the server %+v", server)
			server.GracefulStop()
		})
	}
	for i, server := range servers {
		go serveWithRetries(ctx, server, sm.config.ServerOptions[i], sm.config.ServeRetryDelay, sm.wg)
	}
}

// serveWithRetries tries to serve the server and reties on error.
// Will continue until context is done
// wg done is called when served
func serveWithRetries(ctx context.Context, server *grpc.Server, serverOptions ServerOptions, serverRetryDelay time.Duration, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Debugf("serving with options %+v", serverOptions)
			err := serveServer(server, serverOptions.Port, wg) // blocking
			if err != nil {
				log.Warnf("error serving Server %+v: %+v", serverOptions, err)
			}
		}
		time.Sleep(serverRetryDelay)
	}
}

// serveServer start net.Listener on the port and call done on wg when listening and serves
func serveServer(server *grpc.Server, port int, wg *sync.WaitGroup) error {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	wg.Done()
	if err := server.Serve(list); err != nil {
		return err
	}
	return nil
}
