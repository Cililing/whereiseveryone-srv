package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = time.Duration(30) * time.Second

type Indexable interface {
	EnsureIndexes(ctx context.Context) error
}

type Collections struct {
	Users *mongo.Collection
}

func NewMongoWithPassword(ctx context.Context, db, uri, authDB, user, pass string) (*Collections, error) {
	opts := options.Client().ApplyURI(uri)
	opts.SetServerSelectionTimeout(timeout)
	opts.SetAuth(options.Credential{
		AuthSource: authDB,
		Username:   user,
		Password:   pass,
	})
	return newMongo(ctx, db, opts)
}

func newMongo(ctx context.Context, db string, opts *options.ClientOptions) (*Collections, error) {

	cl, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("create mongo client: %w", err)
	}

	if err := cl.Connect(ctx); err != nil {
		return nil, fmt.Errorf("connect to the db: %w", err)
	}

	if err := cl.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	appDB := cl.Database(db)

	return &Collections{
		Users: appDB.Collection("users"),
	}, nil
}
