package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"whereiseveryone/pkg/id"
	"whereiseveryone/pkg/logger"
	"whereiseveryone/pkg/pointers"
)

type User struct {
	// ID is internal ID
	ID id.ID `bson:"_id"` //nolint:tagliatelle // mongo-id
	// Username is a unique name of player (nick), used for login as well
	Username string `bson:"name"`
	// Password is an encrypted password
	Password string `bson:"password"`
	// Token is a jwt-token
	Token string `bson:"token"`
	// RefreshToken is jwt-refresh-token
	RefreshToken string `bson:"refresh_token"`
	// CreatedAt tells when the user was created
	CreatedAt time.Time `bson:"created_at"`
	// UpdatedAt tells when the last update was done
	UpdatedAt time.Time `bson:"updated_at"`

	// Location user last location (can be nil)
	Location *Location `bson:"location"`
}

type Adapter interface {
	LocationAdapter

	NewUser(ctx context.Context, user User) (User, error)
	GetUserByName(ctx context.Context, name string) (User, error)
	TranslateNamesToIDs(ctx context.Context, names []string) ([]id.ID, error)
	UsersDetails(ctx context.Context, ids []id.ID) (map[id.ID]User, error)
	// UpdateTokens update user tokens (if they are not nil)
	UpdateTokens(ctx context.Context, userID id.ID, token, refreshedToken *string) error
}

var ErrUserNotExists = mongo.ErrNoDocuments
var ErrUserNameAlreadyExists = errors.New("username is already in use")

type mongoAdapter struct {
	LocationAdapter

	coll   *mongo.Collection
	logger logger.Logger
}

func NewMongoAdapter(coll *mongo.Collection, logger logger.Logger) *mongoAdapter {
	locationAdapter := mongoLocationAdapter{coll, logger}
	return &mongoAdapter{locationAdapter, coll, logger}
}

func (m *mongoAdapter) EnsureIndexes(ctx context.Context) error {
	unique := options.IndexOptions{
		Unique: pointers.Pointer(true),
	}
	userIDIdx := mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		},
		Options: &unique,
	}

	_, err := m.coll.Indexes().CreateOne(ctx, userIDIdx)
	if err != nil {
		return fmt.Errorf("create unique name:1 index: %w", err)
	}

	return nil
}

func (m *mongoAdapter) NewUser(ctx context.Context, user User) (User, error) {
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

func (m *mongoAdapter) GetUserByName(ctx context.Context, name string) (User, error) {
	filter := bson.M{
		"name": name,
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

func (m *mongoAdapter) TranslateNamesToIDs(ctx context.Context, names []string) ([]id.ID, error) {
	filter := bson.M{
		"name": bson.M{
			"$in": names,
		},
	}
	projection := bson.M{
		"_id": 1,
	}

	c, err := m.coll.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, fmt.Errorf("perform find query: %w", err)
	}
	defer c.Close(ctx)

	var users []User // will bind only defined projection
	if err := c.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("decode query result with projection: %w", err)
	}

	var res []id.ID
	for _, u := range users {
		res = append(res, u.ID)
	}

	return res, nil
}

func (m *mongoAdapter) UsersDetails(ctx context.Context, ids []id.ID) (map[id.ID]User, error) {
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

	res := make(map[id.ID]User, len(users))
	for _, v := range users {
		res[v.ID] = v
	}
	return res, nil
}

func (m *mongoAdapter) UpdateTokens(ctx context.Context, userID id.ID, token, refreshedToken *string) error {
	tokens := bson.D{}
	if token != nil {
		tokens = append(tokens, bson.E{Key: "token", Value: *token})
	}
	if refreshedToken != nil {
		tokens = append(tokens, bson.E{Key: "refresh_token", Value: *refreshedToken})
	}

	if len(tokens) == 0 {
		// nothing to update
		return nil
	}

	filter := bson.M{
		"_id": userID,
	}
	update := bson.M{
		"$set": tokens,
	}

	_, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update tokens: %w", err)
	}

	return nil
}

var _ Adapter = (*mongoAdapter)(nil)
