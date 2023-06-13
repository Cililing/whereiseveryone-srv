package main

import (
	"context"
	"fmt"
	echoSwagger "github.com/swaggo/echo-swagger"
	"time"

	"github.com/go-playground/validator"
	"whereiseveryone/internal/mongo"
	"whereiseveryone/internal/users"
	"whereiseveryone/internal/webapi"
	authMux "whereiseveryone/internal/webapi/auth"
	locationMux "whereiseveryone/internal/webapi/location"
	"whereiseveryone/pkg/env"
	"whereiseveryone/pkg/jwt"
	"whereiseveryone/pkg/logger"
	"whereiseveryone/pkg/timer"

	_ "github.com/swaggo/echo-swagger" // echo-swagger middleware
	_ "whereiseveryone/docs"
)

// @title WhereIsEveryone
// @version 1.0
// @description This is a sample server for WhereIsEveryone

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api

func main() {
	appCtx := context.Background()
	log := logger.NewLogger()
	utcTimer := timer.NewUTCTimer()

	// Mongo
	mongoURI := env.Env("MONGO_URI", "mongodb://localhost:27017")
	mongoAuthDB := env.Env("MONGO_AUTH_DB", "admin")
	mongoUser := env.Env("MONGO_USER", "root")
	mongoPassword := env.Env("MONGO_PASSWORD", "password123")
	mongoDB := env.Env("MONGO_DB", "whereiseveryone")
	mongoCollections, err := mongo.NewMongo(appCtx, mongoURI, mongoAuthDB, mongoUser, mongoPassword, mongoDB)
	if err != nil {
		log.Fatalf("init mongo: %s", err.Error())
	}
	usersAdapter := users.NewMongoAdapter(mongoCollections.Users, log)

	// Echo
	jwtSecret := env.Env("JWT_SECRET", "jwt-token-123")
	jwtInstance := jwt.NewJWT(utcTimer, []byte(jwtSecret), time.Duration(168)*time.Hour)
	authRouter := authMux.NewMux(usersAdapter, utcTimer, jwtInstance)
	locationRouter := locationMux.NewMux(usersAdapter, log)

	isDebug := env.Env("DEBUG", "true")
	validate := validator.New()
	e := webapi.NewEcho(
		"/api",
		validate,
		jwtInstance,
		webapi.EchoRouters{
			AuthRouter:     authRouter,
			LocationRouter: locationRouter,
		},
		log,
		isDebug == "true")

	// serve docs
	// TODO: move it to webapi????
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	port := env.Env("APP_PORT", "8080")
	log.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
