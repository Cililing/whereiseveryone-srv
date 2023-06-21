package main

import (
	"context"
	"fmt"
	echoSwagger "github.com/swaggo/echo-swagger"
	"time"
	"whereiseveryone/internal/config"

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

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// @BasePath /api

func main() {
	appCtx := context.Background()
	log := logger.NewLogger()
	utcTimer := timer.NewUTCTimer()
	envHandler, err := env.NewHandler("./.env/cloud.json")
	if err != nil {
		log.Fatalf("loading config: %s", err.Error())
	}

	// Mongo
	mongoCollections, err := mongo.GetMongo(appCtx, envHandler)
	if err != nil {
		log.Fatalf("init mongo: %s", err.Error())
	}
	defer mongoCollections.Disconnect(appCtx)
	usersAdapter := users.NewMongoAdapter(mongoCollections.Users, log)

	// Echo
	jwtSecret := envHandler.MustEnv(config.ConfJwtSecret)
	jwtInstance := jwt.NewJWT(utcTimer, []byte(jwtSecret), time.Duration(168)*time.Hour)
	authRouter := authMux.NewMux(usersAdapter, utcTimer, jwtInstance)
	locationRouter := locationMux.NewMux(usersAdapter, usersAdapter, log, utcTimer)

	isDebug := envHandler.MustEnv(config.ConfDebug)
	validate := validator.New()
	e := webapi.NewEcho(
		"/api",
		validate,
		jwtInstance,
		webapi.EchoRouters{
			Swagger:        echoSwagger.WrapHandler,
			AuthRouter:     authRouter,
			LocationRouter: locationRouter,
		},
		log,
		isDebug == "true")

	// Start server
	port := envHandler.MustEnv(config.ConfAppPort)
	log.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
