package location

import "time"

type userLocation struct {
	UUID       string    `json:"id"`
	Nick       string    `json:"nick"`
	Longitude  float64   `json:"longitude"`
	Latitude   float64   `json:"latitude"`
	LastUpdate time.Time `json:"last_update"`
}

type updateLocationRequest struct {
	Longitude float64 `json:"longitude" validate:"required,latitude"`
	Latitude  float64 `json:"latitude" validate:"required,longitude"`
}

// marker, empty response
type updateLocationResponse string

type fetchRequest struct {
	UserIDs []string `json:"uuids"`
	Nicks   []string `json:"nicks"`
}

type fetchUserResponse []userLocation
