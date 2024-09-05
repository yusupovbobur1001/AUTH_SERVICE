package service

import (
	pb "auth_service/genproto/auth_service"
	"auth_service/pkg/logger"
	"auth_service/storage"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceManagement struct {
	storage storage.IUserStorage
	log     logger.ILogger
	*pb.UnimplementedUserServiceServer
}

func (service *UserServiceManagement) Users() storage.IUserStorage {
	return service.storage
}

func NewUserService(storage storage.IUserStorage, log logger.ILogger) *UserServiceManagement {
	return &UserServiceManagement{storage: storage, log: log}
}

func (service *UserServiceManagement) GetUser(ctx context.Context, req *pb.PrimaryKey) (*pb.User, error) {
	return service.storage.GetUser(ctx, req)
}

func (service *UserServiceManagement) GetAllUsers(ctx context.Context, req *pb.GetListRequest) (*pb.GetAllUsersResponse, error) {
	return service.storage.GetAllUsers(ctx, req)
}

func (service *UserServiceManagement) DeleteUser(ctx context.Context, req *pb.PrimaryKey) (*emptypb.Empty, error) {
	return service.storage.DeleteUser(ctx, req)
}

func (service *UserServiceManagement) UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	return service.storage.UpdateUser(ctx, req)
}
