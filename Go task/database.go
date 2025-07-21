package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DBStorage struct {
	DB *sql.DB
}

// add or update scooter information
func (s *DBStorage) AddOrUpdateScooter(ctx context.Context, scooter *Scooter) error {
	_, err := s.DB.ExecContext(ctx, `
        INSERT INTO scooters (id, status, latitude, longitude, last_update)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT(id) DO UPDATE SET
            status=EXCLUDED.status,
            latitude=EXCLUDED.latitude,
            longitude=EXCLUDED.longitude,
            last_update=EXCLUDED.last_update
    `, scooter.ID, scooter.Status, scooter.Latitude, scooter.Longitude, scooter.LastUpdate)
	if err != nil {
		log.Printf("DB error: %v", err)
	}
	return err
}

// function to reserver scooters if they are free (returns true if success)
func (s *DBStorage) ReserveScooter(ctx context.Context, id string) (bool, error) {
	result, err := s.DB.ExecContext(ctx, `
		UPDATE scooters
		SET status = 'occupied'
		WHERE id = $1 AND status = 'free'
	`, id)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (s *DBStorage) FindScooters(ctx context.Context, minLat, minLng, maxLat, maxLng float64, status string) ([]*Scooter, error) {
	query := `SELECT id, status, latitude, longitude, last_update FROM scooters
              WHERE latitude BETWEEN $1 AND $2 AND longitude BETWEEN $3 AND $4`
	args := []interface{}{minLat, maxLat, minLng, maxLng}
	if status != "" {
		query += " AND status=$5"
		args = append(args, status)
	}
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var scooters []*Scooter
	for rows.Next() {
		var sct Scooter
		var ts time.Time
		if err := rows.Scan(&sct.ID, &sct.Status, &sct.Latitude, &sct.Longitude, &ts); err != nil {
			return nil, err
		}
		sct.LastUpdate = ts
		scooters = append(scooters, &sct)
	}
	return scooters, nil
}
