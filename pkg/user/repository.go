package user

import (
	"context"

	pb "github.com/IgorDevCuemby/crudGrpcMysqlMicroservice/users/userpb"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*pb.User, error)
	GetOne(ctx context.Context, id string) (*pb.CreateUserRequest, error)
	Verify(ctx context.Context, id uint32) (*pb.User, error)
	Create(ctx context.Context, user *pb.CreateUserRequest) (*pb.CreateUserRequest, error)
	Update(ctx context.Context, id string, user *pb.CreateUserRequest) (*pb.CreateUserRequest, error)
	Delete(ctx context.Context, id string) error
}
