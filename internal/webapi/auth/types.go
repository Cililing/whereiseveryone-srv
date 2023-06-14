package auth

type signUpRequest struct {
	// Name username, must be unique
	Name string `json:"name" validate:"required"`
	// Password user password, min 8 characters
	Password string `json:"password" validate:"required,min=8"`
	// Email user email (not actually used)
	Email string `json:"email" validate:"required,email"`
}

type logInRequest struct {
	// Name username
	Name string `json:"name" validate:"required"`
	// Password user password
	Password string `json:"password" validate:"required"`
}

type authResponse struct {
	// ID is user id (uuid)
	ID string `json:"id"`
	// Token user auth token (Bearer)
	Token string `json:"token"`
	// RefreshToken user refresh token
	RefreshToken string `json:"refresh_token"`
}
