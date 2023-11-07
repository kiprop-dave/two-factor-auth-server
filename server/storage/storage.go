package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password,omitempty" bson:"password"`
	TagId       string             `json:"tagId" bson:"tagId"`
	TwoFaSecret string             `json:"twoFaSecret,omitempty" bson:"twoFaSecret"`
	TwoFaQrUri  string             `json:"twoFaQrUri,omitempty" bson:"twoFaQrUri"`
	Role        string             `json:"role,omitempty" bson:"role"`
}

type EntryAttempt struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserId     primitive.ObjectID `json:"userId" bson:"userId"`
	TagId      string             `json:"tagId" bson:"tagId"`
	Time       time.Time          `json:"time" bson:"time"`
	Successful bool               `json:"successful" bson:"successful"`
}

type Storage interface {
	GetUser(bson.M) (*User, error)
	CreateUser(user *User) error
	DeleteUser(bson.M) error
	GetUsers() ([]*User, error)
	GetEntryAttempts() ([]*EntryAttempt, error)
	SaveEntryAttempt(entryAttempt *EntryAttempt) error
}

type MongoStorage struct {
	userCollection         *mongo.Collection
	entryAttemptCollection *mongo.Collection
}

func ConnectMongoDb() (*mongo.Database, error) {
	fmt.Println("Connecting to MongoDB...")
	connString := os.Getenv("MONGO_URI")
	if connString == "" {
		return nil, fmt.Errorf("MONGO_URI is not set")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")

	return client.Database("two_fa"), nil
}

func NewMongoStorage(db *mongo.Database) (*MongoStorage, error) {
	userCollection := db.Collection("users")
	entryAttemptCollection := db.Collection("entryAttempts")

	if userCollection == nil {
		return nil, fmt.Errorf("user collection is nil")
	}

	if entryAttemptCollection == nil {
		return nil, fmt.Errorf("entry attempt collection is nil")
	}

	return &MongoStorage{
		userCollection:         userCollection,
		entryAttemptCollection: entryAttemptCollection,
	}, nil
}

func (s *MongoStorage) GetUser(query bson.M) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user := new(User)
	err := s.userCollection.FindOne(ctx, query).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *MongoStorage) CreateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conflictErr := s.userCollection.FindOne(ctx, bson.M{"$or": bson.A{bson.M{"email": user.Email}, bson.M{"tagId": user.TagId}}}).Err()
	if conflictErr == mongo.ErrNoDocuments {
		_, err := s.userCollection.InsertOne(ctx, user)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrConflict
}

func (s *MongoStorage) DeleteUser(bson.M) error {
	return nil
}

func (s *MongoStorage) GetUsers() ([]*User, error) {
	return nil, nil
}

func (s *MongoStorage) GetEntryAttempts() ([]*EntryAttempt, error) {
	return nil, nil
}

func (s *MongoStorage) SaveEntryAttempt(entryAttempt *EntryAttempt) error {
	return nil
}

var ErrConflict = errors.New("conflict")
