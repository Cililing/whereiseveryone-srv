# Run...

## TODO

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
