package servers

// generate a grpc server implementing code in src/generated/userservice
// and run it on port 8080
import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "user-service/src/generated/userservice"
	us "user-service/src/services"
)

// UserServiceServer is the server API for UserService service.
type UserServiceServer interface {
}

type userServiceServerImpl struct {
	pb.UnimplementedUserServiceServer
	userService us.UserService
}

// NewUserServiceServer creates a new UserServiceServer.
func NewUserServiceServer(userService us.UserService) UserServiceServer {
	return &userServiceServerImpl{userService: userService}
}

// CreateUser creates a new user.
func (s *userServiceServerImpl) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return nil, getRequestError(nil)
}

func getRequestError(_ error) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

// GetUser retrieves a user by ID.
func (s *userServiceServerImpl) GetUser(ctx context.Context, req *pb.UserIdRequest) (*pb.UserResponse, error) {
	return nil, getRequestError(nil)
}

// UpdateUser updates an existing user.
func (s *userServiceServerImpl) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return nil, getRequestError(nil)
}

// DeleteUser deletes a user by ID.
func (s *userServiceServerImpl) DeleteUser(ctx context.Context, req *pb.UserIdRequest) (*pb.UserResponse, error) {
	return nil, getRequestError(nil)
}

// ListUsers retrieves all users.
func (s *userServiceServerImpl) ListUsers(ctx context.Context, req *emptypb.Empty) (*pb.UsersResponse, error) {
	return nil, getRequestError(nil)
}

// ResetPassword resets the password for the user
func (s *userServiceServerImpl) ResetPassword(context.Context, *pb.PasswordRequest) (*emptypb.Empty, error) {
	return nil, getRequestError(nil)
}

// VerifyPassword verify the password matches what is in the data store
func (s *userServiceServerImpl) VerifyPassword(context.Context, *pb.PasswordRequest) (*pb.VerifyPasswordResponse, error) {
	return nil, getRequestError(nil)
}
