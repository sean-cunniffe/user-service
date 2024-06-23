package manager

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"user-service/dto"
	pb "user-service/generated/userservice"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// MockUserServiceServer is a mock of UserServiceServer interface
type MockUserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

const (
	ADDR = "localhost"
)

var (
	testTLS = &dto.TLSConfig{
		CAFile:   "../../generated/test_certs/test_ca.crt",
		CertFile: "../../generated/test_certs/test_server.crt",
		KeyFile:  "../../generated/test_certs/test_server.key",
	}
)

func TestStartGrpcServers(t *testing.T) {

	t.Run("test running one nonTLS server", func(t *testing.T) {
		// Create server options
		serverOptions := []ServerOptions{
			{
				Port:              50051,
				ConnectionTimeout: 2 * time.Second,
				TLSConfig:         nil,
				MaxRecvMsgSize:    1024 * 1024,
			},
		}
		testServerWithOptions(t, serverOptions)
	})

	t.Run("test running one TLS server", func(t *testing.T) {
		// Create server options
		serverOptions := []ServerOptions{
			{
				Port:              50052,
				ConnectionTimeout: 2 * time.Second,
				TLSConfig:         testTLS,
				MaxRecvMsgSize:    1024 * 1024,
			},
		}
		testServerWithOptions(t, serverOptions)
	})
	t.Run("test running one TLS server and one non TLS server", func(t *testing.T) {
		// Create server options
		serverOptions := []ServerOptions{
			{
				Port:              50052,
				ConnectionTimeout: 2 * time.Second,
				TLSConfig:         testTLS,
				MaxRecvMsgSize:    1024 * 1024,
			},
			{
				Port:              50051,
				ConnectionTimeout: 2 * time.Second,
				TLSConfig:         nil,
				MaxRecvMsgSize:    1024 * 1024,
			},
		}
		testServerWithOptions(t, serverOptions)
	})
	t.Run("test servers never ready if server cannot start", func(t *testing.T) {
		// create two servers on same port so second one cannot start
		serverOptions := []ServerOptions{
			{
				Port:              50051,
				ConnectionTimeout: 2 * time.Second,
				TLSConfig:         nil,
				MaxRecvMsgSize:    1024 * 1024,
			},
			{
				Port:              50051,
				ConnectionTimeout: 2 * time.Second,
				TLSConfig:         nil,
				MaxRecvMsgSize:    1024 * 1024,
			},
		}
		testContext, cancelCtx := context.WithCancel(context.Background())
		ctrl := gomock.NewController(t)
		mockUserServiceServer := &MockUserServiceServer{}
		defer ctrl.Finish()
		waitServerStop := time.Millisecond * 2000
		// Create server manager
		sm := CreateServerManager(ServerManagerConfig{ServeRetryDelay: time.Millisecond * 500, ServerOptions: serverOptions}, mockUserServiceServer)
		// Start servers
		sm.StartGrpcServers(testContext)
		go func(ctx context.Context) {
			select {
			case <-ctx.Done():
				return
			case <-sm.ServersStarted():
				assert.FailNow(t, "servers should not be ready")

			}
		}(testContext)
		// wait for server to stop
		time.Sleep(waitServerStop)
		cancelCtx()
	})
}

func testServerWithOptions(t *testing.T, serverOptions []ServerOptions) {
	testContext, cancelCtx := context.WithCancel(context.Background())
	ctrl := gomock.NewController(t)
	mockUserServiceServer := &MockUserServiceServer{}
	defer ctrl.Finish()
	waitServerStop := time.Millisecond * 500
	// Create server manager
	sm := CreateServerManager(ServerManagerConfig{ServeRetryDelay: time.Millisecond * 1000, ServerOptions: serverOptions}, mockUserServiceServer)
	// Start servers
	sm.StartGrpcServers(testContext)
	<-sm.ServersStarted()
	for _, options := range serverOptions {
		assert.True(t, isPortOpen(t, options.Port, 2*time.Second))
	}

	// Stop servers
	cancelCtx()
	// wait for server to stop
	time.Sleep(waitServerStop)
	// Check if server is still open after stop signal
	for _, options := range serverOptions {
		assert.False(t, isPortOpen(t, options.Port, 2*time.Second))
	}
}

func isPortOpen(t *testing.T, port int, timeout time.Duration) bool {
	t.Helper()
	address := fmt.Sprintf("%s:%d", ADDR, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
