package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"userManagement/infra/database"
	pb "userManagement/proto"
)

var (
	client = database.DBClient
)

func TestDBCreateUser(t *testing.T) {
	createdID, err := client.CreateUser(&pb.CreateUserReq{User: testUser})
	if err != nil {
		return
	}

	assert.EqualValues(t, 24, len(createdID))
}

func TestDBCreateRepeatedUser(t *testing.T) {
	_, err := client.CreateUser(&pb.CreateUserReq{User: testUser})
	if err == nil {
		t.Fatal("Create a user with an already registered email should not be permitted")
	}
}

func TestDBGetUser(t *testing.T) {
	user, err := client.GetUser(&pb.GetUserReq{UserId: userID})
	if err != nil {
		t.Fatal("Could not retrieve user")
	}

	assert.EqualValues(t, testResponse.User.Email, user.User.Email)
}

func TestDBGetAllUsersFilterLess(t *testing.T) {
	users, err := client.GetAllUsers(&pb.ListUsersReq{Filter: nil})

	if err != nil {
		t.Fatalf("Failed when retrieving users")
	}

	assert.EqualValues(t, "testing", users.Users[0].User.FirstName)
}

func TestDBGetAllUsersWithFilter(t *testing.T) {
	users, err := client.GetAllUsers(&pb.ListUsersReq{
		Filter: &pb.User{
			Country: "ES",
		},
	})

	if err != nil {
		t.Fatalf("Failed when retrieving users")
	}

	assert.EqualValues(t, "testing", users.Users[0].User.FirstName)
}

func TestDBUpdateUser(t *testing.T) {
	user, err := client.UpdateUser(&pb.UpdateUserReq{
		UserId: userID,
		User: &pb.User{
			FirstName: "testing-updated",
			LastName:  "user",
			Email:     "a@a.com",
			Nickname:  "a",
			Password:  "1234",
			Country:   "ES",
		},
	})

	if err != nil {
		t.Fatal("Could not update user")
	}

	assert.EqualValues(t, user.User.FirstName, "testing-updated")
}

func TestDBDeleteUser(t *testing.T) {
	_, err := client.DeleteUser(&pb.DeleteUserReq{UserId: userID})
	if err != nil {
		t.Fatal("Could not delete user")
	}
}

func TestDBGetMissingUser(t *testing.T) {
	_, err := client.GetUser(&pb.GetUserReq{UserId: userID})
	if err == nil {
		t.Fatal("This test is supposed to retrieve a non-existing user, an error should occur")
	}
}
