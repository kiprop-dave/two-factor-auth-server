package storage

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserId primitive.ObjectID `json:"userId" bson:"userId"`
	//CheckPointId primitive.ObjectID `json:"checkPointId" bson:"checkPointId"`
	TagId      string    `json:"tagId" bson:"tagId"`
	Time       time.Time `json:"time" bson:"time"`
	Successful bool      `json:"successful" bson:"successful"`
}

type CheckPoint struct {
	Name   string             `json:"name" bson:"name"`
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	ApiKey string             `json:"apiKey" bson:"apiKey"`
}

type Storage interface {
	GetUser(bson.M) (*User, error)
	CreateUser(user *User) error
	DeleteUser(bson.M) error
	GetUsers() ([]*User, error)
	GetEntryAttempts() ([]*EntryAttempt, error)
	SaveEntryAttempt(entryAttempt *EntryAttempt) (string, error)
	SaveSuccessfulEntryAttempt(entryAttemptId string) error
	GetEntryAttempt(id string) (*EntryAttempt, error)
	CreateCheckPoint(name string) (*CheckPoint, error)
}

type MongoStorage struct {
	userCollection         *mongo.Collection
	entryAttemptCollection *mongo.Collection
	checkPointCollection   *mongo.Collection
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
	checkPointCollection := db.Collection("checkPoints")

	if userCollection == nil {
		return nil, fmt.Errorf("user collection is nil")
	}

	if entryAttemptCollection == nil {
		return nil, fmt.Errorf("entry attempt collection is nil")
	}

	if checkPointCollection == nil {
		return nil, fmt.Errorf("check point collection is nil")
	}

	return &MongoStorage{
		userCollection:         userCollection,
		entryAttemptCollection: entryAttemptCollection,
		checkPointCollection:   checkPointCollection,
	}, nil
}

func (s *MongoStorage) GetUser(query bson.M) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user := new(User)
	err := s.userCollection.FindOne(ctx, query).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
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

func (s *MongoStorage) SaveEntryAttempt(entryAttempt *EntryAttempt) (string, error) {
	res, err := s.entryAttemptCollection.InsertOne(context.Background(), entryAttempt)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *MongoStorage) GetEntryAttempt(id string) (*EntryAttempt, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	attempt := new(EntryAttempt)
	if err := s.entryAttemptCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(attempt); err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("not found")
			return nil, ErrNotFound
		}
		return nil, err
	}
	return attempt, nil
}

func (s *MongoStorage) SaveSuccessfulEntryAttempt(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.entryAttemptCollection.UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": bson.M{"successful": true}})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoStorage) CreateCheckPoint(name string) (*CheckPoint, error) {
	buff := make([]byte, 16)
	_, err := rand.Read(buff)
	if err != nil {
		return nil, err
	}

	newCheckPoint := CheckPoint{
		Name:   name,
		ID:     primitive.NewObjectID(),
		ApiKey: base64.StdEncoding.EncodeToString(buff),
	}
	_, err = s.checkPointCollection.InsertOne(context.Background(), newCheckPoint)

	if err != nil {
		return nil, err
	}

	return &newCheckPoint, nil
}

var ErrConflict = errors.New("conflict")
var ErrNotFound = errors.New("not found")
