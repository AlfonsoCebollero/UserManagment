package server

import (
	"context"
	"log"
	"time"
	"userManagement/infra/database"
	pb "userManagement/proto"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	DbClient database.AdapterInterface
}

func (s *UserManagementServer) CreateUser(ctx context.Context, in *pb.CreateUserReq) (*pb.UserActionResponse, error) {
	user := in.GetUser()
	log.Printf("Received: %v", user)

	// Store new user in database
	userId, err := s.DbClient.CreateUser(in)
	if err != nil {
		return nil, err
	}

	return &pb.UserActionResponse{
		Id: userId,
		User: &pb.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Password:  user.Password,
			Country:   user.Country,
		},
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}, nil
}
