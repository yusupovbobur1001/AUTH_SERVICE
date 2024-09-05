package storage

import (
	pb "auth_service/genproto/auth_service"
	"auth_service/models"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

type IUserStorage interface {
	GetUser(ctx context.Context, request *pb.PrimaryKey) (*pb.User, error)
	GetAllUsers(ctx context.Context, request *pb.GetListRequest) (*pb.GetAllUsersResponse, error)
	DeleteUser(ctx context.Context, request *pb.PrimaryKey) (*empty.Empty, error)
	UpdateUser(ctx context.Context, request *pb.User) (*pb.User, error)
}

type IAuthStorage interface {
	RegisterRepo(ctx context.Context, request *models.CreateUser) (*models.User, error)
	Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error)
	ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) (*models.Message, error)
	ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) (*models.Message, error)
}
