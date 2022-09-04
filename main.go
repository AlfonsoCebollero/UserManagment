package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"userManagement/infra/database"
	"userManagement/infra/server"
	pb "userManagement/proto"
)

const (
	port = ":5566"
)

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "swagger/swagger.json")
}

func runAPIServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Connect to the GRPC server
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register grpc-gateway
	rmux := runtime.NewServeMux()
	err := pb.RegisterUserManagementHandlerFromEndpoint(ctx, rmux, "localhost:5566", dopts)
	if err != nil {
		log.Fatal(err)
	}

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.HandleFunc("/swagger.json", serveSwagger)
	sh := http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger/")))
	mux.Handle("/swagger/", sh)
	log.Println("REST server ready...")
	err = http.ListenAndServe("localhost:8081", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	go runAPIServer()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, &server.UserManagementServer{
		DbClient:      database.DBClient,
		NotifyChannel: make(chan []string),
	})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}
