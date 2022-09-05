package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "userManagement/proto"
)

const (
	address = "host.docker.internal:5566"
)

func main() {
	log.Print("Connecting to server")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	log.Printf("Connected to grpc server: %v", conn)
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Could not close connection %v", err)
		}
	}(conn)

	c := pb.NewUserManagementClient(conn)

	changeStream, err := c.NotifyUserChanges(context.TODO(), &pb.EmptyMsg{})
	log.Println("Server side streaming established")
	if err != nil {
		log.Printf("There was an error when recieving changes: %v", err)
		return
	}

	for {
		log.Print("Waiting for new notifications")
		var notification pb.UserActionStream
		err := changeStream.RecvMsg(&notification)
		if err != nil {
			log.Printf("Failed when recieving change: %v", err)
		}
		log.Printf(notification.Action)
	}
}
