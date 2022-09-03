package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/mail"
	"time"
	"userManagement/entities"
	pb "userManagement/proto"
)

var (
	DBClient *MongoClient
)

func init() {
	// mongo connection initialization
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// generation of a unique client to interact with mongo
	DBClient = &MongoClient{
		Collection: client.
			Database("userManagement").
			Collection("users")}
	log.Printf("%v", DBClient)
}

type MongoClient struct {
	Collection *mongo.Collection
}

// CreateUser adds a new user to the database.
func (m *MongoClient) CreateUser(req *pb.CreateUserReq) (string, error) {
	user := req.GetUser()

	_, err := mail.ParseAddress(user.Email)

	if err != nil {
		return "", entities.InvalidEmailError
	}

	filter := bson.D{primitive.E{Key: "email", Value: user.Email}}

	var foundUser bson.M
	err = m.Collection.FindOne(context.Background(), filter).Decode(&foundUser)
	if err == mongo.ErrNoDocuments {
		log.Printf("User email is not registered, registering new user...")
	} else {
		if err == nil {
			log.Printf("Coulf not create user: %v", entities.AlreadyRegisteredEmailError)
			return "", entities.AlreadyRegisteredEmailError
		}
		log.Printf("Could not create user: %v", err)
		return "", err
	}

	mongoUser := entities.User{
		Id:        primitive.NewObjectID(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.GetNickname(),
		Password:  user.Password,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := m.Collection.InsertOne(context.Background(), mongoUser)
	if err != nil {
		log.Printf("Could not create user with mail %s, %v", user.Email, err)
		return "", err
	}

	createdID := createdUser.InsertedID.(primitive.ObjectID).Hex()
	return createdID, nil
}

// GetUser retrieves a user from the database
func (m *MongoClient) GetUser(req *pb.GetUserReq) (*pb.UserActionResponse, error) {
	id := req.UserId
	filter := getFindUserFilter(id)

	var foundUser entities.User
	err := m.Collection.FindOne(context.Background(), filter).Decode(&foundUser)
	if err != nil {

		msg := "Could not find user with id %s: %v"
		return nil, handleActionError(id, msg, err)
	}

	return getPbUser(foundUser)

}

// UpdateUser finds a user inside the database and update its fields.
// Email cannot be updated since is used along _id to identify unique users
func (m *MongoClient) UpdateUser(req *pb.UpdateUserReq) (*pb.UserActionResponse, error) {
	id := req.UserId
	filter := getFindUserFilter(id)
	user := req.User

	update := bson.D{{"$set", bson.D{
		{"first_name", user.FirstName},
		{"last_name", user.LastName},
		{"nickname", user.Nickname},
		{"password", user.Password},
		{"country", user.Country},
		{"updated_at", time.Now()}}}}

	var updatedUser entities.User
	err := m.Collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&updatedUser)

	if err != nil {
		msg := "Could not update user with id %s: %v"
		return nil, handleActionError(id, msg, err)
	}

	req.ProtoMessage()

	return getPbUser(updatedUser)
}

// DeleteUser removes a user from the database
func (m *MongoClient) DeleteUser(req *pb.DeleteUserReq) (*pb.DeletionActionResponse, error) {
	id := req.UserId

	filter := getFindUserFilter(id)

	_, err := m.Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		msg := "Could not update user with id %s: %v"
		return nil, handleActionError(id, msg, err)
	}

	return nil, nil
}

// GetAllUsers Filters by the provided query params in the request and returns all users that match the filter
func (m *MongoClient) GetAllUsers(req *pb.ListUsersReq) (*pb.ListActionResponse, error) {
	filter := req.Filter

	filterMap := entities.User{
		FirstName: filter.FirstName,
		LastName:  filter.LastName,
		Nickname:  filter.Nickname,
		Email:     filter.Email,
		Country:   filter.Country,
	}
	var mongoFilter bson.D
	mFilter, err := bson.Marshal(filterMap)
	if err != nil {
		return nil, err
	}
	_ = bson.Unmarshal(mFilter, &mongoFilter)

	cursor, err := m.Collection.Find(context.TODO(), mongoFilter)
	if err != nil {
		return nil, err
	}
	var results []entities.User
	pbResults := pb.ListActionResponse{
		Users: []*pb.UserActionResponse{},
	}

	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		user, err := getPbUser(result)
		if err != nil {
			return nil, err
		}
		pbResults.Users = append(pbResults.Users, user)
	}

	return &pbResults, nil
}

func handleActionError(id, msg string, err error) error {
	log.Printf(msg, id, err)
	if err == mongo.ErrNoDocuments {
		return entities.NotFoundUser
	}
	return err
}

// getFindUserFilter is an auxiliar function that builds the filter to find users.
// It is necessary to be able to filter by email and by mongo id
func getFindUserFilter(id string) bson.D {
	_, err := mail.ParseAddress(id)
	var filter bson.D
	if err == nil {
		filter = bson.D{primitive.E{Key: "email", Value: id}}
	} else {
		mongoID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.D{primitive.E{Key: "_id", Value: mongoID}}
	}
	return filter
}

// getPbUser builds a protobuf user response from the user entity used to interact with mongo
func getPbUser(foundUser entities.User) (*pb.UserActionResponse, error) {
	return &pb.UserActionResponse{
		Id: foundUser.Id.Hex(),
		User: &pb.User{
			FirstName: foundUser.FirstName,
			LastName:  foundUser.LastName,
			Email:     foundUser.Email,
			Nickname:  foundUser.Nickname,
			Password:  foundUser.Password,
			Country:   foundUser.Country,
		},
		CreatedAt: foundUser.CreatedAt.String(),
		UpdatedAt: foundUser.UpdatedAt.String(),
	}, nil
}
