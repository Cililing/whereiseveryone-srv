package users

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"slices"
	"whereiseveryone/pkg/id"
	"whereiseveryone/pkg/logger"
	"whereiseveryone/pkg/pointers"
	"whereiseveryone/pkg/timer"
)

type User struct {
	// ID is internal ID
	ID id.ID `bson:"_id"` //nolint:tagliatelle // mongo-id
	// Auth is auth details for the user (JWT)
	Auth Auth `bson:"auth"`
	// Location user last location (can be nil)
	Location *Location `bson:"location"`

	// User text status (can be empty)
	Status string `bson:"status"`

	// ObservedUsers list of IDs user subscribe
	// NOTE: The second user should accept subscription.
	//		 For now, before returning those user data,
	//		 we need to make sure both users subscribes each other.
	SubscribedUsers []id.ID `bson:"subscribed_users"`
}

func (u User) SubscribeUser(id id.ID) bool {
	return slices.Contains(u.SubscribedUsers, id)
}

type Adapter interface {
	locationAdapter
	authAdapter

	NewUser(ctx context.Context, user User) (User, error)

	GetUser(ctx context.Context, userID id.ID) (User, error)
	GetUsers(ctx context.Context, ids []id.ID) ([]User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)

	UpdateStatus(ctx context.Context, user id.ID, newStatus string) error
	ObserveUser(ctx context.Context, user id.ID, userToObserve id.ID) error
	UnobserveUser(ctx context.Context, user id.ID, userToUnobserve id.ID) error
}

var ErrUserNotExists = mongo.ErrNoDocuments
var ErrUserNameAlreadyExists = errors.New("username is already in use")

type mongoUserAdapter struct {
	locationAdapter
	authAdapter

	coll   *mongo.Collection
	logger logger.Logger
}

func NewMongoAdapter(coll *mongo.Collection, timer timer.Timer, logger logger.Logger) *mongoUserAdapter {
	locationAdapter := mongoLocationAdapter{coll, logger}
	authAdapter := mongoAuthAdapter{coll, timer, logger}

	return &mongoUserAdapter{locationAdapter, authAdapter, coll, logger}
}

func (m *mongoUserAdapter) EnsureIndexes(ctx context.Context) error {
	unique := options.IndexOptions{
		Unique: pointers.Pointer(true),
	}
	userIDIdx := mongo.IndexModel{
		Keys: bson.M{
			"auth.name": 1,
		},
		Options: &unique,
	}

	_, err := m.coll.Indexes().CreateOne(ctx, userIDIdx)
	if err != nil {
		return fmt.Errorf("create unique name:1 index: %w", err)
	}

	return nil
}

func (m *mongoUserAdapter) NewUser(ctx context.Context, user User) (User, error) {
	user.ID = id.NewID()
	_, err := m.coll.InsertOne(ctx, user)
	if err != nil {
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, innerErr := range writeErr.WriteErrors {
				if innerErr.Code == 11000 { // duplicate err
					return User{}, ErrUserNameAlreadyExists
				}
			}
		}
		return User{}, fmt.Errorf("create a new user: %w", err)
	}

	return user, nil
}

func (m *mongoUserAdapter) GetUser(ctx context.Context, userID id.ID) (User, error) {
	filter := bson.M{
		"_id": userID,
	}

	res := m.coll.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return User{}, fmt.Errorf("find user by id: %w", err)
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (m *mongoUserAdapter) GetUsers(ctx context.Context, ids []id.ID) ([]User, error) {
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	c, err := m.coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("perform find query: %w", err)
	}

	var users []User
	if err := c.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("decode query result: %w", err)
	}

	return users, nil
}

func (m *mongoUserAdapter) GetUserByUsername(ctx context.Context, name string) (User, error) {
	filter := bson.M{
		"auth.name": name,
	}

	res := m.coll.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return User{}, fmt.Errorf("perform query: %w", err)
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return User{}, fmt.Errorf("decode query result: %w", err)
	}

	return user, nil
}

func (m *mongoUserAdapter) UpdateStatus(ctx context.Context, userId id.ID, newStatus string) error {
	filter := bson.M{
		"id": userId,
	}
	update := bson.M{
		"$set": bson.D{
			{Key: "status", Value: newStatus},
		},
	}

	_, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update user's status: %w", err)
	}

	return nil
}

func (m *mongoUserAdapter) ObserveUser(ctx context.Context, user id.ID, userToObserve id.ID) error {
	return nil
}

func (m *mongoUserAdapter) UnobserveUser(ctx context.Context, user id.ID, userToUnobserve id.ID) error {
	return nil
}

var _ Adapter = (*mongoUserAdapter)(nil)
