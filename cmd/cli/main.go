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
	envHandler, err := env.NewHandler("./.env/cloud.json")
	if err != nil {
		log.Fatalf("loading config: %s", err.Error())
	}

	mongoCollections, err := mongo.GetMongo(appCtx, envHandler)
	if err != nil {
		log.Fatalf("connecting with mongo: %s", err.Error())
	}

	cmd := commands.NewCommandApp(log, utcTimer, mongoCollections)
	if err := cmd.ExecuteContext(appCtx); err != nil {
		log.Fatalf("executing command: %s", err.Error())
	}
}
