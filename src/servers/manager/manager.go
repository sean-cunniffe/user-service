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
	"google.golang.org/grpc/reflection"
)

type (
	// ServerManager manages the lifecycle of the grpc servers.
	// ServerManager exposes functions to start the servers based on
	// its config
	ServerManager interface {
		StartGrpcServers(ctx context.Context)
		ServersStarted() <-chan struct{}
		OnError(func(err error))
		OnServing(func())
	}

	// ServerManagerConfig holds the config for the servers it manages
	ServerManagerConfig struct {
		ServeRetryDelay time.Duration   `yaml:"serveRetryDelay" json:"serveRetryDelay"`
		ServerOptions   []ServerOptions `yaml:"serverOptions" json:"serverOptions"`
	}
	// ServerOptions contains the properties passed into the grpc server
	ServerOptions struct {
		Port              int            `yaml:"port" json:"port"`
		ConnectionTimeout time.Duration  `yaml:"connectionTimeout" json:"connectionTimeout"`
		TLSConfig         *dto.TLSConfig `yaml:"tlsConfig" json:"tlsConfig"`
		MaxRecvMsgSize    int            `yaml:"maxRecvMsgSize" json:"maxRecvMsgSize"`
	}
	serverManager struct {
		config            ServerManagerConfig
		userServiceServer pb.UserServiceServer
		wg                *sync.WaitGroup
		onErrorFunc       func(err error)
		onServingFunc     func()
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
		onErrorFunc: func(err error) {
			log.Warnf("error in server manager: %+v", err)
		},
		onServingFunc: func() {
			log.Info("server is serving")
		},
	}
}

func (sm *serverManager) OnError(f func(err error)) {
	sm.onErrorFunc = f
}

func (sm *serverManager) OnServing(f func()) {
	sm.onServingFunc = f
}

func getGrpcServerOptions(serverOptions ServerOptions) ([]grpc.ServerOption, error) {
	options := []grpc.ServerOption{
		grpc.ConnectionTimeout(serverOptions.ConnectionTimeout),
		grpc.MaxRecvMsgSize(serverOptions.MaxRecvMsgSize),
	}
	if serverOptions.TLSConfig == nil {
		log.Printf("running server without TLS")
		return options, nil
	}
	// Load TLS config
	tlsConfig, err := util.GetTLSConfig(*serverOptions.TLSConfig)
	if err != nil {
		return nil, err
	}
	options = append(options, grpc.Creds(credentials.NewTLS(tlsConfig)))
	return options, nil
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
		grpcServerOptions, err := getGrpcServerOptions(options)
		if err != nil {
			sm.onErrorFunc(err)
			return
		}
		server := grpc.NewServer(grpcServerOptions...)
		pb.RegisterUserServiceServer(server, sm.userServiceServer)
		reflection.Register(server)
		servers[i] = server
		context.AfterFunc(ctx, func() {
			log.Debugf("shutting down the server %+v", server)
			server.GracefulStop()
		})
	}
	for i, server := range servers {
		go serveWithRetries(ctx, server, sm.config.ServerOptions[i], sm.config.ServeRetryDelay, sm.wg, sm.onErrorFunc, sm.onServingFunc)
	}
}

// serveWithRetries tries to serve the server and reties on error.
// Will continue until context is done
// wg done is called when served
func serveWithRetries(ctx context.Context, server *grpc.Server, serverOptions ServerOptions, serverRetryDelay time.Duration, wg *sync.WaitGroup, onErr func(err error), onServing func()) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Debugf("serving with options %+v", serverOptions)
			err := serveServer(server, serverOptions.Port, wg, onServing) // blocking
			if err != nil {
				log.Warnf("error serving Server %+v: %+v", serverOptions, err)
				onErr(err)
			}
		}
		time.Sleep(serverRetryDelay)
	}
}

// serveServer start net.Listener on the port and call done on wg when listening and serves
func serveServer(server *grpc.Server, port int, wg *sync.WaitGroup, onServing func()) error {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	wg.Done()
	onServing()
	if err := server.Serve(list); err != nil {
		return err
	}
	return nil
}
