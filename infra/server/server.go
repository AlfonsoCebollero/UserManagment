package server

import (
	"context"
	"fmt"
	"log"
	"userManagement/infra/database"
	pb "userManagement/proto"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	DbClient      database.AdapterInterface
	NotifyChannel chan []string
}

// CreateUser creates a new user from the received request and returns user details
// It sends a creation action notification
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

	go s.notify(in.User.String(), "Created")
	return createdUser, nil
}

// GetUser retrieves a user by id or by email
// it sends a retrieving action notification
func (s *UserManagementServer) GetUser(ctx context.Context, in *pb.GetUserReq) (*pb.UserActionResponse, error) {
	user, err := s.DbClient.GetUser(in)
	if err != nil {
		return nil, err
	}

	go s.notify(in.UserId, "Retrieved")
	return user, nil
}

// UpdateUser updates a user, who is found by email or ID. It uses the body of the request to update.
// It sends an update action notification.
func (s *UserManagementServer) UpdateUser(ctx context.Context, in *pb.UpdateUserReq) (*pb.UserActionResponse, error) {
	user, err := s.DbClient.UpdateUser(in)
	if err != nil {
		return nil, err
	}

	go s.notify(in.UserId, "Updated")
	return user, nil
}

// DeleteUser removes a user from the database by ID or email
// Sends a deletion action notification
func (s *UserManagementServer) DeleteUser(ctx context.Context, in *pb.DeleteUserReq) (*pb.DeletionActionResponse, error) {
	_, err := s.DbClient.DeleteUser(in)

	if err != nil {
		return &pb.DeletionActionResponse{Deleted: false}, err
	}

	go s.notify(in.UserId, "Deleted")
	return &pb.DeletionActionResponse{Deleted: true}, nil
}

// ListUsers retrieves all stored users
func (s *UserManagementServer) ListUsers(ctx context.Context, in *pb.ListUsersReq) (*pb.ListActionResponse, error) {
	users, err := s.DbClient.GetAllUsers(in)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// NotifyUserChanges creates a stream where action notifications are received.
func (s *UserManagementServer) NotifyUserChanges(msg *pb.EmptyMsg, server pb.UserManagement_NotifyUserChangesServer) error {
	for {
		select {
		case n := <-s.NotifyChannel:
			notification := pb.UserActionStream{}
			action := fmt.Sprintf("User action performed: %s - %s", n[0], n[1])
			notification.Action = action

			err := server.Send(&notification)
			if err != nil {
				return err
			}
		}
	}
}

// notify sends a notification through the server channel to later be sent through
// a server side streaming
func (s *UserManagementServer) notify(userEmail, action string) {
	notification := []string{userEmail, action}
	s.NotifyChannel <- notification
	return
}
