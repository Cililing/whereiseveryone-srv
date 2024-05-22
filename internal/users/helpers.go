package users

import (
	"go.mongodb.org/mongo-driver/bson"
	"whereiseveryone/pkg/id"
)

func withUserId(id id.ID) bson.M {
	return bson.M{
		"_id": id,
	}
}
