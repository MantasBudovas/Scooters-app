package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Launches fake client
func SimulateClient(clientID string) {
	client := &http.Client{}

	for {
		// 1. Checks for free scooters
		scooters, err := getScootersWithAuth(client, "http://localhost:8080/scooters?min_lat=45.0&max_lat=46.0&min_lng=-76.0&max_lng=-73.0&status=free")
		if err != nil {
			log.Println("Failed to get scooters:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if len(scooters) == 0 {
			log.Println("No free scooters, waiting...")
			time.Sleep(2 * time.Second)
			continue
		}
		// Chooses a random free scooter
		s := scooters[rand.Intn(len(scooters))]

		// 2. Try start trip
		startEvent := Event{
			ID:        uuid.New().String(),
			ScooterID: s.ID,
			Type:      "start",
			Latitude:  s.Latitude,
			Longitude: s.Longitude,
			Timestamp: time.Now(),
		}

		err = sendEvent(client, startEvent)
		if err != nil {
			log.Printf("[%s] failed to start ride on scooter %s: %v", clientID, s.ID, err)
			time.Sleep(1 * time.Second)
			continue
		}
		log.Printf("[%s] started ride on scooter %s", clientID, s.ID)

		// 3. Traveling for 10-15s + location updates every 3s
		duration := rand.Intn(6) + 10
		steps := duration / 3
		latStep := 0.001
		lngStep := 0.001
		for i := 0; i < steps; i++ {
			s.Latitude += latStep
			s.Longitude += lngStep
			updateEvent := Event{
				ID:        uuid.New().String(),
				ScooterID: s.ID,
				Type:      "update",
				Latitude:  s.Latitude,
				Longitude: s.Longitude,
				Timestamp: time.Now(),
				Status:    "occupied",
			}
			sendEvent(client, updateEvent)
			log.Printf("[%s] %s - update %d: location [lat: %.5f, lng: %.5f]",
				clientID, s.ID, i+1, s.Latitude, s.Longitude)
			time.Sleep(3 * time.Second)
		}

		// 4. Ends trip
		endEvent := Event{
			ID:        uuid.New().String(),
			ScooterID: s.ID,
			Type:      "end",
			Latitude:  s.Latitude,
			Longitude: s.Longitude,
			Timestamp: time.Now(),
		}
		sendEvent(client, endEvent)
		log.Printf("[%s] finished ride on scooter %s, ride time: %d seconds", clientID, s.ID, duration)

		// Rest 2-5 s
		time.Sleep(time.Duration(rand.Intn(4)+2) * time.Second)
	}
}

// Get scooters with API key
func getScootersWithAuth(client *http.Client, url string) ([]Scooter, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-API-Key", ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var scooters []Scooter
	if err := json.NewDecoder(resp.Body).Decode(&scooters); err != nil {
		return nil, err
	}
	return scooters, nil
}

func sendEvent(client *http.Client, ev Event) error {
	data, _ := json.Marshal(ev)
	req, _ := http.NewRequest("POST", "http://localhost:8080/event", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}
	return nil
}

// Simulates x clients
func StartSimulation() {
	for i := 1; i <= 3; i++ {
		go SimulateClient("client-" + strconv.Itoa(i))
	}
}
