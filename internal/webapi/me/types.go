package me

import "time"

type updateStatusRequest struct {
	Status string `json:"status"`
}

type getFriendsResponse []friendDetails

type friendDetails struct {
	Username string          `json:"username"`
	Status   string          `json:"status"`
	Location locationDetails `json:"location"`
}

type locationDetails struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Altitude  float64 `json:"altitude,omitempty"`
	Bearing   float64 `json:"bearing,omitempty"`
	Accuracy  float64 `json:"accuracy,omitempty"`
	// LastUpdate in UTC time
	LastUpdate time.Time `json:"last_update"`
}

type updateLocationRequest struct {
	locationDetails `json:",inline"`
}

type observeRequest struct {
	Username string `json:"username"`
}
