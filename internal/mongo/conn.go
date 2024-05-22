package mongo

import (
	"context"
	"fmt"
	"strconv"
	"whereiseveryone/internal/config"
	"whereiseveryone/pkg/env"
)

func GetMongo(appCtx context.Context, env env.Handler) (*Collections, error) {
	mongoUseCloud := env.MustEnv(config.ConfMongoUseCloud)
	useCloud, _ := strconv.ParseBool(mongoUseCloud)

	if useCloud {
		return CloudMongo(appCtx, env)
	}

	return LocalMongo(appCtx, env)
}

func LocalMongo(appCtx context.Context, env env.Handler) (*Collections, error) {
	mongoURI := env.MustEnv(config.ConfMongoUri)
	mongoAuthDB := env.MustEnv(config.ConfMongoAuthDb)
	mongoUser := env.MustEnv(config.ConfMongoUser)
	mongoPassword := env.MustEnv(config.ConfMongoPassword)
	mongoDB := env.MustEnv(config.ConfMongoDb)

	mongoCollections, err := NewMongoWithPassword(appCtx, mongoDB, mongoURI, mongoAuthDB, mongoUser, mongoPassword)
	if err != nil {
		return nil, fmt.Errorf("init mongo: %w", err)
	}
	return mongoCollections, nil
}

func CloudMongo(appCtx context.Context, env env.Handler) (*Collections, error) {
	mongoURI := env.MustEnv(config.ConfMongoUri)
	mongoCrt := env.MustEnv(config.ConfMongoX509)
	mongoDB := env.MustEnv(config.ConfMongoDb)

	mongoCollections, err := NewMongoWithX509Pem(appCtx, mongoDB, mongoURI, mongoCrt)
	if err != nil {
		return nil, fmt.Errorf("init mongo: %w", err)
	}
	return mongoCollections, nil
}
