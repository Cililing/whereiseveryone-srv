package users

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"whereiseveryone/pkg/id"
	"whereiseveryone/pkg/logger"
)

type Location struct {
	// Longitude
	Longitude float64 `bson:"longitude"`
	// Latitude
	Latitude float64 `bson:"latitude"`
	// Altitude
	Altitude float64 `bson:"altitude,omitempty"`
	// Bearing
	Bearing float64 `bson:"bearing,omitempty"`
	// Accuracy
	Accuracy float64 `bson:"accuracy,omitempty"`
	// LastUpdate
	LastUpdate time.Time `bson:"last_update"`
}

type locationAdapter interface {
	UpdateLocation(ctx context.Context, userID id.ID, newLocation Location) error
}

type mongoLocationAdapter struct {
	userCollection *mongo.Collection

	logger logger.Logger
}

func (l mongoLocationAdapter) UpdateLocation(ctx context.Context, userID id.ID, newLocation Location) error {
	location := bson.D{
		bson.E{Key: "location", Value: newLocation},
	}

	filter := withUserId(userID)
	update := bson.M{
		"$set": location,
	}

	_, err := l.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update user's location: %w", err)
	}

	return nil
}

var _ locationAdapter = (*mongoLocationAdapter)(nil)
