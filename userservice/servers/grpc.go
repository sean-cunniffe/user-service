package servers

// generate a grpc server implementing code in generated/userservice
// and run it on port 8080
import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "user-service/generated/userservice"
	us "user-service/services"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService us.UserService
}

// CreateUserServiceServer creates the grpc server instance which is ran on a grpc server
func CreateUserServiceServer(userService us.UserService) pb.UserServiceServer {
	return &userServiceServer{
		userService: userService,
	}
}

// NewUserServiceServer creates a new UserServiceServer instance.

// CreateUser creates a new user.
func (s *userServiceServer) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log.Infof("received request %+v", req)
	return nil, getRequestError(nil)
}

func getRequestError(_ error) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

// GetUser retrieves a user by ID.
func (s *userServiceServer) GetUser(ctx context.Context, req *pb.UserIdRequest) (*pb.UserResponse, error) {
	log.Infof("received request %+v", req)
	return nil, getRequestError(nil)
}

// UpdateUser updates an existing user.
func (s *userServiceServer) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log.Infof("received request %+v", req)
	return nil, getRequestError(nil)
}

// DeleteUser deletes a user by ID.
func (s *userServiceServer) DeleteUser(ctx context.Context, req *pb.UserIdRequest) (*pb.UserResponse, error) {
	log.Infof("received request %+v", req)
	return nil, getRequestError(nil)
}

// ListUsers retrieves all users.
func (s *userServiceServer) ListUsers(ctx context.Context, req *emptypb.Empty) (*pb.UsersResponse, error) {
	log.Infof("received request %+v", req)
	return nil, getRequestError(nil)
}

// ResetPassword resets the password for the user
func (s *userServiceServer) ResetPassword(ctx context.Context, req *pb.PasswordRequest) (*emptypb.Empty, error) {
	log.Infof("received request %+v", req)
	return nil, getRequestError(nil)
}

// VerifyPassword verify the password matches what is in the data store
func (s *userServiceServer) VerifyPassword(ctx context.Context, req *pb.PasswordRequest) (*pb.VerifyPasswordResponse, error) {
	return nil, getRequestError(nil)
}
