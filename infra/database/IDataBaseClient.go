package database

import pb "userManagement/proto"

type AdapterInterface interface {
	CreateUser(req *pb.CreateUserReq) (string, error)
	GetUser(req *pb.GetUserReq) (*pb.UserActionResponse, error)
	UpdateUser(req *pb.UpdateUserReq) (*pb.UserActionResponse, error)
	DeleteUser(req *pb.DeleteUserReq) (*pb.DeletionActionResponse, error)
	GetAllUsers(req *pb.ListUsersReq) (*pb.ListActionResponse, error)
}
