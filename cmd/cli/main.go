package main

import (
	"context"
	"flag"
	"whereiseveryone/cmd/cli/commands"
	"whereiseveryone/internal/mongo"
	"whereiseveryone/pkg/env"
	"whereiseveryone/pkg/logger"
	"whereiseveryone/pkg/timer"
)

func main() {
	configPathFlag := flag.String("config", "./.env/local.json", "config path")
	flag.Parse()
	envHandler, err := env.NewHandler(*configPathFlag)

	appCtx := context.Background()
	log := logger.NewLogger()
	utcTimer := timer.NewUTCTimer()
	if err != nil {
		log.Fatalf("loading config: %s", err.Error())
	}

	mongoCollections, err := mongo.GetMongo(appCtx, envHandler)
	if err != nil {
		log.Fatalf("connecting with mongo: %s", err.Error())
	}
	defer mongoCollections.Disconnect(appCtx)

	cmd := commands.NewCommandApp(log, utcTimer, mongoCollections)
	if err := cmd.ExecuteContext(appCtx); err != nil {
		log.Fatalf("executing command: %s", err.Error())
	}
}
