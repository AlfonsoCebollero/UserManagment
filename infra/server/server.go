package server

import (
	"context"
	"log"
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
	log.Printf("User created!")

	createdUser, err := s.GetUser(ctx, &pb.GetUserReq{UserId: userId})

	if err != nil {
		log.Printf("Could not retrieve created user data")
		return nil, err
	}

	return createdUser, nil
}

func (s *UserManagementServer) GetUser(ctx context.Context, in *pb.GetUserReq) (*pb.UserActionResponse, error) {
	user, err := s.DbClient.GetUser(in)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserManagementServer) UpdateUser(ctx context.Context, in *pb.UpdateUserReq) (*pb.UserActionResponse, error) {
	user, err := s.DbClient.UpdateUser(in)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserManagementServer) DeleteUser(ctx context.Context, in *pb.DeleteUserReq) (*pb.DeletionActionResponse, error) {
	_, err := s.DbClient.DeleteUser(in)

	if err != nil {
		return &pb.DeletionActionResponse{Deleted: false}, err
	}

	return &pb.DeletionActionResponse{Deleted: true}, nil
}

func (s *UserManagementServer) ListUsers(ctx context.Context, in *pb.ListUsersReq) (*pb.ListActionResponse, error) {
	users, err := s.DbClient.GetAllUsers(in)
	if err != nil {
		return nil, err
	}

	return users, nil
}
