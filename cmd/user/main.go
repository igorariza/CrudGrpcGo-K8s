package main

import (
	"os"

	"github.com/IgorDevCuemby/crudGrpcMysqlMicroservice/internal/data"
	"github.com/IgorDevCuemby/crudGrpcMysqlMicroservice/pkg/user"

	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/IgorDevCuemby/crudGrpcMysqlMicroservice/users/userpb"
	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	Repository user.Repository
}

var users []*pb.User

//Listar
func (s *server) ListUser(ctx context.Context, req *empty.Empty) (*pb.ListUsersResponse, error) {
	users_all, err := s.Repository.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return &pb.ListUsersResponse{
		Users: users_all,
	}, nil
}

//Create user
func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserId, error) {
	log.Printf("Creating User %v", req)
	user, err := s.Repository.Create(ctx, req)

	if err != nil {
		return nil, err
	}
	return &pb.UserId{
		Id: user.Id,
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.UserId) (*pb.Response, error) {
	err := s.Repository.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Msg: fmt.Sprintf("DELETE %v", req.Id)}, nil
}

func (s *server) GetOneArticle(ctx context.Context, req *pb.UserId) (*pb.CreateUserRequest, error) {
	log.Printf("Get User By ID %v", req.Id)

	article, err := s.Repository.GetOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (s *server) UpdateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserRequest, error) {

	user, err := s.Repository.Update(ctx, req.Id, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *server) VerifyUser(ctx context.Context, req *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {
	log.Printf("Get User If Exist %v", req.UserId)

	_, err := s.Repository.Verify(ctx, req.UserId)
	if err != nil {
		return &pb.VerifyUserResponse{IsExist: 0}, nil
	}
	return &pb.VerifyUserResponse{IsExist: 1}, nil
}

//
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	fmt.Println("User")

	if err != nil {
		log.Fatalf("Error cannot create tcp connection %v", err)
	}

	d := data.New()
	if err := d.DB.Ping(); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatalf("Error cannot create tcp connection %v", err)
	}
	log.Printf("Connection established running on port %v", port)

	ser := grpc.NewServer()
	pb.RegisterUserServiceServer(ser, &server{
		Repository: &data.UsersRepository{
			Data: data.New(),
		},
	})
	if err := ser.Serve(listen); err != nil {
		log.Fatalf("Error cannot initialize the server: %v", err.Error())
	}

	defer data.Close()
}
