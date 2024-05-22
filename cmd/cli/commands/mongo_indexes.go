package commands

import (
	"whereiseveryone/internal/users"
)

func (c *commandApp) mongoIndexes() {
	usersAdapter := users.NewMongoAdapter(c.mongoColls.Users, c.logger)
	if err := usersAdapter.EnsureIndexes(c.Context()); err != nil {
		c.logger.Fatalf("create indexes on users collection: %s", err.Error())
	}
}
