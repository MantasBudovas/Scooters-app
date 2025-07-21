package main

import "time"

type Scooter struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"` // "free" or "occupied"
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	LastUpdate time.Time `json:"last_update"`
}

type Event struct {
	ID        string    `json:"id"`
	ScooterID string    `json:"scooter_id"`
	Type      string    `json:"type"` // "start", "end", "update"
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"` // "free" or "occupied"
}
