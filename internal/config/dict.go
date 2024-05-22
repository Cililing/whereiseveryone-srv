package config

import "whereiseveryone/pkg/env"

const (
	ConfMongoUseCloud env.Key = "mongo.useCloud" // required
	ConfMongoUri              = "mongo.uri"      // required
	ConfMongoAuthDb           = "mongo.authDb"   // required for local-db
	ConfMongoUser             = "mongo.user"     // required for local-db
	ConfMongoPassword         = "mongo.password" // required for local-db
	ConfMongoDb               = "mongo.db"       // required
	ConfMongoX509             = "mongo.x509"     // required for cloud

	ConfJwtSecret = "app.jwtSecret" // required
	ConfDebug     = "app.debug"     // required
	ConfAppPort   = "app.port"      // required
)
