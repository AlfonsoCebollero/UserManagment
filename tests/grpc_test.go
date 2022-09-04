package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"
	"userManagement/infra/server"
	pb "userManagement/proto"
)

const (
	bufSize = 1024 * 1024
	userID  = "a@a.com"
)

var (
	lis *bufconn.Listener

	testResponse = &pb.UserActionResponse{
		Id: "1",
		User: &pb.User{
			FirstName: "testing",
			LastName:  "user",
			Email:     "a@a.com",
			Nickname:  "a",
			Password:  "1234",
			Country:   "ES",
		},
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	testUser = &pb.User{
		FirstName: "testing",
		LastName:  "user",
		Email:     "a@a.com",
		Nickname:  "a",
		Password:  "1234",
		Country:   "ES",
	}

	grpcServer = server.UserManagementServer{
		DbClient:      nil,
		NotifyChannel: make(chan []string),
	}
)

type DBAdapterMock struct {
	mock.Mock
}

func (m *DBAdapterMock) GetUser(req *pb.GetUserReq) (*pb.UserActionResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*pb.UserActionResponse), args.Error(1)
}

func (m *DBAdapterMock) UpdateUser(req *pb.UpdateUserReq) (*pb.UserActionResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*pb.UserActionResponse), args.Error(1)
}

func (m *DBAdapterMock) DeleteUser(req *pb.DeleteUserReq) (*pb.DeletionActionResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*pb.DeletionActionResponse), args.Error(1)
}

func (m *DBAdapterMock) GetAllUsers(req *pb.ListUsersReq) (*pb.ListActionResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*pb.ListActionResponse), args.Error(1)
}

func (m *DBAdapterMock) CreateUser(req *pb.CreateUserReq) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func TestGetUser(t *testing.T) {
	mockDBClient := new(DBAdapterMock)

	grpcServer.DbClient = mockDBClient

	mockDBClient.On("GetUser", &pb.GetUserReq{UserId: userID}).Return(testResponse, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	resp, err := grpcServer.GetUser(ctx, &pb.GetUserReq{UserId: userID})
	if err != nil {
		t.Fatalf("Get User test failed: %v", err)
	}

	mockDBClient.AssertExpectations(t)

	assert.EqualValues(t, testResponse, resp)
	return
}

func TestUpdateUser(t *testing.T) {
	mockDBClient := new(DBAdapterMock)

	grpcServer.DbClient = mockDBClient

	mockDBClient.On("UpdateUser", &pb.UpdateUserReq{
		UserId: userID,
		User:   testUser,
	}).Return(testResponse, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	resp, err := grpcServer.UpdateUser(ctx, &pb.UpdateUserReq{
		UserId: userID,
		User:   testUser,
	})
	if err != nil {
		t.Fatalf("Get User test failed: %v", err)
	}

	mockDBClient.AssertExpectations(t)

	assert.EqualValues(t, testResponse, resp)
}

func TestDeleteUser(t *testing.T) {
	mockDBClient := new(DBAdapterMock)

	grpcServer.DbClient = mockDBClient

	mockDBClient.On("DeleteUser", &pb.DeleteUserReq{
		UserId: userID,
	}).Return(&pb.DeletionActionResponse{Deleted: true}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	resp, err := grpcServer.DeleteUser(ctx, &pb.DeleteUserReq{
		UserId: userID,
	})
	if err != nil {
		t.Fatalf("Get User test failed: %v", err)
	}

	mockDBClient.AssertExpectations(t)

	assert.EqualValues(t, true, resp.Deleted)
}

func TestListUsers(t *testing.T) {
	mockDBClient := new(DBAdapterMock)

	grpcServer.DbClient = mockDBClient

	mockDBClient.On("GetAllUsers", &pb.ListUsersReq{}).
		Return(&pb.ListActionResponse{
			Users: []*pb.UserActionResponse{testResponse},
		}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	resp, err := grpcServer.ListUsers(ctx, &pb.ListUsersReq{Filter: nil})
	if err != nil {
		t.Fatalf("Get User test failed: %v", err)
	}

	mockDBClient.AssertExpectations(t)

	assert.EqualValues(t, testResponse, resp.Users[0])
}

func TestCreateUser(t *testing.T) {
	mockDBClient := new(DBAdapterMock)
	grpcServer.DbClient = mockDBClient

	mockDBClient.On("CreateUser", &pb.CreateUserReq{
		User: testUser,
	}).Return(userID, nil)

	mockDBClient.On("GetUser", &pb.GetUserReq{UserId: userID}).
		Return(testResponse, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	resp, err := grpcServer.CreateUser(ctx, &pb.CreateUserReq{User: testUser})
	if err != nil {
		t.Fatalf("Get User test failed: %v", err)
	}

	mockDBClient.AssertExpectations(t)

	assert.EqualValues(t, resp, testResponse)
}

// TestNotifyChanges establishes a server side streaming to receive notifications
// from user changes.
func TestNotifyChanges(t *testing.T) {
	mockDBClient := new(DBAdapterMock)
	grpcServer.DbClient = mockDBClient

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	mockDBClient.On("GetUser", &pb.GetUserReq{UserId: userID}).
		Return(testResponse, nil)

	pb.RegisterUserManagementServer(s, &grpcServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.
		DialContext(ctx,
			"bufnet",
			grpc.WithContextDialer(bufDialer),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserManagementClient(conn)
	resp, err := client.NotifyUserChanges(ctx, &pb.EmptyMsg{})
	if err != nil {
		t.Fatalf("Error when creating server side stream: %v", err)
	}

	_, err = grpcServer.GetUser(ctx, &pb.GetUserReq{UserId: userID})
	if err != nil {
		return
	}

	var notification pb.UserActionStream
	err = resp.RecvMsg(&notification)
	if err != nil {
		t.Fatalf("Error when recieving user action notification")
	}

	assert.EqualValues(t, "User action performed: a@a.com - Retrieved", notification.Action)

	mockDBClient.AssertExpectations(t)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
