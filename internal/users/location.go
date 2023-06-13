package users

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"whereiseveryone/pkg/id"
	"whereiseveryone/pkg/logger"
)

type Location struct {
	// Longitude
	Longitude float64 `bson:"longitude"`
	// Latitude
	Latitude float64 `bson:"latitude"`
	// LastUpdate
	LastUpdate time.Time `bson:"last_update"`
}

type LocationAdapter interface {
	UpdateLocation(ctx context.Context, userID id.ID, newLocation Location) error
	GetLocation(ctx context.Context, userID id.ID) (*Location, error)
}

type mongoLocationAdapter struct {
	userCollection *mongo.Collection

	logger.Logger
}

func (l mongoLocationAdapter) UpdateLocation(ctx context.Context, userID id.ID, newLocation Location) error {
	location := bson.D{
		bson.E{Key: "location", Value: newLocation},
	}

	filter := bson.M{
		"_id": userID,
	}
	update := bson.M{
		"$set": location,
	}

	_, err := l.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update user's location: %w", err)
	}

	return nil
}

func (l mongoLocationAdapter) GetLocation(ctx context.Context, userID id.ID) (*Location, error) {
	filter := bson.M{
		"_id": userID,
	}

	res := l.userCollection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("perform query: %w", err)
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return nil, fmt.Errorf("decode query result: %w", err)
	}

	return user.Location, nil
}

var _ LocationAdapter = (*mongoLocationAdapter)(nil)
