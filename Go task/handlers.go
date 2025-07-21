package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Function for incoming scooter events (start, update, end)
func EventHandlerDB(storage *DBStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		switch event.Type {
		case "start":
			// Reserve the scooter atomically
			reserved, err := storage.ReserveScooter(r.Context(), event.ScooterID)
			if err != nil {
				http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if !reserved {
				http.Error(w, "Scooter is not free", http.StatusConflict)
				return
			}

			// Update location
			scooter := &Scooter{
				ID:         event.ScooterID,
				Latitude:   event.Latitude,
				Longitude:  event.Longitude,
				LastUpdate: event.Timestamp,
				Status:     "occupied",
			}
			err = storage.AddOrUpdateScooter(r.Context(), scooter)
			if err != nil {
				http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
				return
			}

		case "update":
			// Updates location
			scooter := &Scooter{
				ID:         event.ScooterID,
				Latitude:   event.Latitude,
				Longitude:  event.Longitude,
				LastUpdate: event.Timestamp,
				Status:     "occupied",
			}
			err := storage.AddOrUpdateScooter(r.Context(), scooter)
			if err != nil {
				http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
				return
			}

		case "end":
			// Set status to free and update location
			scooter := &Scooter{
				ID:         event.ScooterID,
				Latitude:   event.Latitude,
				Longitude:  event.Longitude,
				LastUpdate: event.Timestamp,
				Status:     "free",
			}
			err := storage.AddOrUpdateScooter(r.Context(), scooter)
			if err != nil {
				http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Unknown event type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"ok"}`))
	}
}

// Handle scooter location/status queries
func ScootersHandler(storage *DBStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minLat, _ := strconv.ParseFloat(r.URL.Query().Get("min_lat"), 64)
		maxLat, _ := strconv.ParseFloat(r.URL.Query().Get("max_lat"), 64)
		minLng, _ := strconv.ParseFloat(r.URL.Query().Get("min_lng"), 64)
		maxLng, _ := strconv.ParseFloat(r.URL.Query().Get("max_lng"), 64)
		status := r.URL.Query().Get("status")
		scooters, err := storage.FindScooters(r.Context(), minLat, minLng, maxLat, maxLng, status)
		if err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(scooters)
	}
}
