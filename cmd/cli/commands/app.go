package commands

import (
	"github.com/spf13/cobra"
	"whereiseveryone/internal/mongo"
	"whereiseveryone/pkg/logger"
	"whereiseveryone/pkg/timer"
)

type commandApp struct {
	*cobra.Command

	logger logger.Logger
	timer  timer.Timer

	mongoColls *mongo.Collections
}

func NewCommandApp(
	lg logger.Logger,
	tm timer.Timer,
	mongoColls *mongo.Collections,
) *commandApp {
	rootCmd := &cobra.Command{
		Use:   "wie", // where is everyone abbrv
		Short: "Some useful commands for where is everyone server app",
		Long: `
			The app allows to execute some useful commands,
  			required to manage the app
			Can be easily used with k8s-jobs or any other scheduler.
			`,
	}

	app := &commandApp{
		Command:    rootCmd,
		logger:     lg,
		timer:      tm,
		mongoColls: mongoColls,
	}

	dummyCmd := &cobra.Command{
		Use:   "dummy",
		Short: "Just print debug info on the screen",
		Run: func(cmd *cobra.Command, args []string) {
			app.logger.Info("You are a dummy!")
		},
	}

	mongoIndexes := &cobra.Command{
		Use:   "mongoIndexes",
		Short: "sets mongo indexes (must be run after each index change)",
		Run: func(cmd *cobra.Command, args []string) {
			app.mongoIndexes()
		},
	}

	rootCmd.AddCommand(dummyCmd)
	rootCmd.AddCommand(mongoIndexes)

	return app
}
