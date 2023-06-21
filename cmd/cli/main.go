package main

import (
	"context"
	"whereiseveryone/cmd/cli/commands"
	"whereiseveryone/pkg/logger"
	"whereiseveryone/pkg/timer"
)

func main() {
	appCtx := context.Background()
	log := logger.NewLogger()
	utcTimer := timer.NewUTCTimer()

	cmd := commands.NewCommandApp(log, utcTimer)
	if err := cmd.ExecuteContext(appCtx); err != nil {
		log.Fatalf("executing command: %s", err.Error())
	}
}
