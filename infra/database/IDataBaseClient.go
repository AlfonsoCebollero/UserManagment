package database

import pb "userManagement/proto"

type AdapterInterface interface {
	CreateUser(req *pb.CreateUserReq) (string, error)
	//UpdateUser(proto.UpdateUserReq) error
	//DeleteUser(proto.DeleteUserReq) error
	//GetAllUsers() error
}
