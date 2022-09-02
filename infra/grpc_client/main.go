package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	pb "userManagement/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Could not close connection %v", err)
		}
	}(conn)

	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, err := c.CreateUser(ctx, &pb.CreateUserReq{
		User: &pb.User{
			FirstName: "Alfonso",
			LastName:  "Cebollero",
			Email:     "a@gmail.com",
			Nickname:  nil,
			Password:  "1234",
			Country:   "ES",
		},
	})
	if err != nil {
		log.Fatalf("Error while creating user %v", err)
	}
	log.Printf("user details %v", user)
	return

}
