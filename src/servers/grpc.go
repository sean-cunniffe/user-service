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

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService us.UserService
}

func CreateUserServiceServer(userService us.UserService) pb.UserServiceServer {
	return &userServiceServer{
		userService: userService,
	}
}

// NewUserServiceServer creates a new UserServiceServer instance.

// CreateUser creates a new user.
func (s *userServiceServer) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return nil, getRequestError(nil)
}

func getRequestError(_ error) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

// GetUser retrieves a user by ID.
func (s *userServiceServer) GetUser(ctx context.Context, req *pb.UserIdRequest) (*pb.UserResponse, error) {
	return nil, getRequestError(nil)
}

// UpdateUser updates an existing user.
func (s *userServiceServer) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return nil, getRequestError(nil)
}

// DeleteUser deletes a user by ID.
func (s *userServiceServer) DeleteUser(ctx context.Context, req *pb.UserIdRequest) (*pb.UserResponse, error) {
	return nil, getRequestError(nil)
}

// ListUsers retrieves all users.
func (s *userServiceServer) ListUsers(ctx context.Context, req *emptypb.Empty) (*pb.UsersResponse, error) {
	return nil, getRequestError(nil)
}

// ResetPassword resets the password for the user
func (s *userServiceServer) ResetPassword(context.Context, *pb.PasswordRequest) (*emptypb.Empty, error) {
	return nil, getRequestError(nil)
}

// VerifyPassword verify the password matches what is in the data store
func (s *userServiceServer) VerifyPassword(context.Context, *pb.PasswordRequest) (*pb.VerifyPasswordResponse, error) {
	return nil, getRequestError(nil)
}
