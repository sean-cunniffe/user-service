package manager

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"user-service/dto"
	pb "user-service/generated/userservice"
	util "user-service/util"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type (
	// ServerManager manages the lifecycle of the grpc servers.
	// ServerManager exposes functions to start the servers based on
	// its config
	// ServerManager also exposes functions to register callbacks
	ServerManager interface {
		StartManager(ctx context.Context)
		OnServing(func())
		OnStopServing(func(serverOptions ServerOptions, err error))
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
		serverErrChan     chan struct {
			so  *ServerOptions
			err error
		}
		serverStartChan   chan *ServerOptions
		onServingFunc     func()
		onStopServingFunc func(ServerOptions, error)
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
		serverErrChan: make(chan struct {
			so  *ServerOptions
			err error
		}),
		serverStartChan: make(chan *ServerOptions),
		onStopServingFunc: func(so ServerOptions, err error) {
			log.Warnf("error in server manager: %+v", err)
		},
		onServingFunc: func() {
			log.Info("server is serving ")
		},
	}
}

func (sm *serverManager) OnStopServing(f func(ServerOptions, error)) {
	sm.onStopServingFunc = f
}

func (sm *serverManager) OnServing(f func()) {
	sm.onServingFunc = f
}

// start the grpc server from the serveroptions
func (sm *serverManager) StartManager(ctx context.Context) {
	servers := make([]*grpc.Server, len(sm.config.ServerOptions))
	serversReady := make(map[*ServerOptions]bool, len(servers))
	for i := 0; i < len(sm.config.ServerOptions); i++ {
		options := &sm.config.ServerOptions[i]
		serversReady[options] = false
		select {
		case <-ctx.Done():
			return
		default:
			server, err := sm.registerServer(ctx, options)
			if err != nil {
				log.Warnf("error registering server %+v: %+v", options, err)
				time.Sleep(sm.config.ServeRetryDelay)
				i = i - 1
				continue
			}
			servers[i] = server
		}
	}
	// start serving
	for i, server := range servers {
		serverOptions := &sm.config.ServerOptions[i]
		go sm.serveWithRetries(ctx, server, serverOptions)
	}
	// watch servers ready
	go sm.watchServersReady(ctx, serversReady)
}

func (sm *serverManager) getGrpcServerOptions(ctx context.Context, serverOptions *ServerOptions) ([]grpc.ServerOption, error) {
	options := []grpc.ServerOption{
		grpc.ConnectionTimeout(serverOptions.ConnectionTimeout),
		grpc.MaxRecvMsgSize(serverOptions.MaxRecvMsgSize),
	}
	if serverOptions.TLSConfig == nil {
		log.Printf("running server without TLS")
		return options, nil
	}
	// Load TLS config
	tlsConfig, err := util.GetRefreshableTLSConfig(ctx, *serverOptions.TLSConfig, func(err error) {
		log.Errorf("failed to refresh TLS configuration: %v", err)
		sm.serverErrChan <- struct {
			so  *ServerOptions
			err error
		}{so: serverOptions, err: err}
	}, func() {
		log.Info("refreshed TLS configuration")
		sm.serverStartChan <- serverOptions
	})
	if err != nil {
		return nil, err
	}
	options = append(options, grpc.Creds(credentials.NewTLS(tlsConfig)))
	return options, nil
}

func (sm *serverManager) watchServersReady(ctx context.Context, serversReady map[*ServerOptions]bool) {
	for {
		select {
		case serverOptions := <-sm.serverStartChan:
			serversReady[serverOptions] = true
			// check nothing in map is false
			allServing := true
			for _, serving := range serversReady {
				if !serving {
					allServing = false
					break
				}
			}
			if allServing {
				sm.onServingFunc()
			}
		case sigerr := <-sm.serverErrChan:
			serversReady[sigerr.so] = false
			sm.onStopServingFunc(*sigerr.so, sigerr.err)
		case <-ctx.Done():
			return
		}
	}
}

func (sm *serverManager) registerServer(ctx context.Context, options *ServerOptions) (*grpc.Server, error) {
	grpcServerOptions, err := sm.getGrpcServerOptions(ctx, options)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer(grpcServerOptions...)
	pb.RegisterUserServiceServer(server, sm.userServiceServer)
	reflection.Register(server)
	return server, nil
}

// serveWithRetries tries to serve the server and reties on error.
// Will continue until context is done
// wg done is called when served
func (sm *serverManager) serveWithRetries(ctx context.Context, server *grpc.Server, serverOptions *ServerOptions) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Debugf("serving with options %+v", serverOptions)
			err := serveServer(ctx, server, serverOptions.Port, func() { sm.serverStartChan <- serverOptions }) // blocking
			if err != nil {
				sm.serverErrChan <- struct {
					so  *ServerOptions
					err error
				}{so: serverOptions, err: err}
			}
		}
		time.Sleep(sm.config.ServeRetryDelay)
	}
}

// serveServer serves the server on the port and calls the onServing function
// returns whem context is done or server is stopped
func serveServer(ctx context.Context, server *grpc.Server, port int, onServing func()) error {
	serverContext, cancel := context.WithCancel(ctx)
	defer cancel()
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	onServing()
	var serverErr error
	go func() {
		if err := server.Serve(list); err != nil {
			serverErr = err
			cancel()
		}
	}()
	<-serverContext.Done()
	server.Stop()
	return serverErr
}
