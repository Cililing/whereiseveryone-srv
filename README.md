# Run

## Go

It's a standard go app. You can run it using `go run` etc.

## Config

App uses json-config. The config MUST be a JSON with **only string** entries.
As a default `./.env/local.json` is used. You can override it with a flag `--config=$filePath`

## Docker

To build a image locally run: `docker build -t whereiseveryone-srv:latest -f docker/Dockerfile .` in project root.
The command will build the container and will append local `./.env` directory. It can be overwritten later.

If you want to use local-running mongo (see section `Development/Mongo` below) in docker a network must be created at
first.

```
docker network create whereiseveryone-net
docker network connect whereiseveryone-net mongodb
```

Then, to run an image:

```
docker run \
    -p 127.0.0.1:8080:8080/tcp \
    --network=whereiseveryone-net \
    -v "`pwd`/.env:/app/.env" \
    whereiseveryone-srv \
    /bin/sh -c "/app/app-srv -config=/app/.env/docker.json"
```

The command:

* `-p` binds docker to your localhost on port 8080
* `--network` connects the container to the network (required for connecting with local-docker mongo)
* `-v` binds local `./env` directory to container `.env` directory
* `/bin/sh -c ...` command for running the srv with docker-config file

# Authorization

- All users are required to create an account.
- For authentication a JWT token is required.
- To signup use `/auth/signup`

```go
package auth

type signUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type logInRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type authResponse struct {
	ID           string `json:"id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
```

All requests (except `/auth/*`) are required to have JWT token attached (`Header -> Authorization: Bearer <<token>>`).
When token expires a user must renew it (with `/login`).

# Development

To run app in development, at first run MongoDB docker container:
`docker run --name mongodb -p 27017:27017 -e MONGODB_ROOT_PASSWORD=password123 bitnami/mongodb:4.4`

This command will run the mongodb container with root user: `root:password123` on port 27017

## Config

To see a list of available config keys please see `/internal/config/dict.go`.
For local development use `./.env/local.json` file.

_To replace a config you need to edit main.go files to point a proper one_

## Using cloud db

At first, you need to generate a X509 certificate from Mongo Atlas. **Keep it secret!**
Put it in `.env` directory and then use a config from `./.env/cloud.json`

# Documentation

Docs are served in /swagger endpoint.
Ref: https://github.com/swaggo/echo-swagger

For generating docs (required each time something is changed) `swag init -g cmd/server/main.go`
and commit it to the repository.

## Binding Requests

There is a very useful generic function that binds the HTTP request and validates it.

```go
package request

func echoFunc(c echo.Context) error {
	data, bindErr := binder.BindRequest[bodyType](c, true)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	return c.String(200, "ok")
}
```

`BindRequest` returns an object implementing the interface

```go
package request

type BaseContext interface {
	Context() context.Context
	Cancel() context.CancelFunc
	Echo() echo.Context
	UserID() id.ID
	TokenData() jwt.SignedToken
}
```

# Production

TBD.
