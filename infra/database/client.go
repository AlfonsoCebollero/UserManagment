package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"userManagement/entities"
	pb "userManagement/proto"
)

var (
	DBClient *MongoClient
)

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	DBClient = &MongoClient{
		Collection: client.
			Database("userManagement").
			Collection("users")}
	log.Printf("%v", DBClient)
}

type MongoClient struct {
	Collection *mongo.Collection
}

func (m *MongoClient) CreateUser(req *pb.CreateUserReq) (string, error) {
	user := req.GetUser()

	filter := bson.D{primitive.E{Key: "email", Value: user.Email}}

	var foundUser bson.M
	err := m.Collection.FindOne(context.Background(), filter).Decode(&foundUser)
	if err == mongo.ErrNoDocuments {
		log.Printf("User email is not registered, registering new user...")
	} else {
		if err == nil {
			log.Printf("Coulf not create user: %v", entities.AlreadyRegisteredEmailError)
			return "", entities.AlreadyRegisteredEmailError
		}
		log.Printf("Coulf not create user: %v", err)
		return "", err
	}

	mongoUser := entities.User{
		Id:        primitive.ObjectID{},
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
		log.Printf("Could not create user %v", err)
		return "", err
	}

	createdID := createdUser.InsertedID.(primitive.ObjectID).Hex()
	return createdID, nil
}
