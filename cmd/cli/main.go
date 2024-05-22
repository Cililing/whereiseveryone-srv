package main

import (
	"context"

	"whereiseveryone/cmd/cli/commands"
	"whereiseveryone/internal/mongo"
	"whereiseveryone/pkg/env"
	"whereiseveryone/pkg/logger"
	"whereiseveryone/pkg/timer"
)

func main() {
	appCtx := context.Background()
	log := logger.NewLogger()
	utcTimer := timer.NewUTCTimer()

	mongoURI := env.Env("MONGO_URI", "mongodb://localhost:27017")
	mongoAuthDB := env.Env("MONGO_AUTH_DB", "admin")
	mongoUser := env.Env("MONGO_USER", "root")
	mongoPassword := env.Env("MONGO_PASSWORD", "password123")
	mongoDB := env.Env("MONGO_DB", "whereiseveryone")
	mongoCollections, err := mongo.NewMongoWithPassword(appCtx, mongoDB, mongoURI, mongoAuthDB, mongoUser, mongoPassword)
	if err != nil {
		log.Fatalf("init mongo: %s", err.Error())
	}

	cmd := commands.NewCommandApp(log, utcTimer, mongoCollections)
	if err := cmd.ExecuteContext(appCtx); err != nil {
		log.Fatal(err)
	}
}
