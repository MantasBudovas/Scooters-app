# Scootin' Aboot — Scooter Event Collection Backend

  

This project is a backend service for managing electric scooters in Ottawa and Montreal. It supports tracking scooter trips and locations and exposes a simple API for mobile clients to interact with.

  ## Table of Contents
- [Features](#Features)
- [Technologies Used](#technologies-used)
- [Setup](#setup)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)
- [Client Simulator](#client-simulator)
- [Database](#database)
- [Testing](#testing)

---

### Features
  

- REST-like API for event reporting and querying scooter statuses.
- Clients simulate scooter use (start, update, end) with geolocation.
- Atomic reservation to avoid double-booking scooters.
- PostgreSQL for persistent data storage.
- Docker Compose for quick deployment.

---
### Technologies Used

- Go (Golang)

- PostgreSQL

- Docker + Docker Compose

- UUIDs for ID management
---

### Setup

### 1. Clone the repository

```bash
git  clone  https://github.com/NordSecurity-Interviews/BE-MantasBudovas.git
```
```
cd  BE-MantasBudovas
```

### 2. Build the application

```bash
docker-compose  up  --build
```

---
### Authentication

All requests require a Static API key

```bash
  X-API-Key: APIkey
```
---
### API Endpoints

1. Post /event
Used by client simulator to report ```start```, ```update``` or ```end``` events.
```
{
	id: Event ID (string)
	scooter_id: Scooter UUID (string)
	type: "start"|"update"|"end" (string)
	latitude: Latitude of the scooter (float)
	longitude: Longitude of the scooter (float)
	timestamp: ISO 8601 timestamp (string)
	status: "occupied"|"free" (string)
}
```
Example request:
```
Content-Type: application/json
X-API-Key: APIkey

{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "scooter_id": "11111111-c073-4e7e-814e-5d6737f4696e",
  "type": "start",
  "latitude": 45.4215,
  "longitude": -75.6972,
  "timestamp": "2025-04-30T12:00:00Z",
  "status": "occupied"
}
```

2. Get /scooters
Searches for scooters by filtering a rectangular location
```
{
	min_lat: Minimum latitude (flaot)
    max_lat: Maximum latitude (float)
    min_lng: Minimum longitude (float)
    max_lng: Maximum longitude (flaot)
    status: "occupied|free" (string)
}
```
Example request:
```
GET /scooters?min_lat=45.0&max_lat=46.0&min_lng=-76.0&max_lng=-73.0&status=occupied
```
---
### Client Simulator
- Searches for free scooters
- Starts trips
- Sends location updates every 3 seconds
- Ends trips and repeat after resting
---
### Database
- PostgreSQL is used as the backing data store
- Schema is defined in ```schema.sql``` and loaded via Docker.
- The ```scooters``` table uses ```UUID``` as the primary identifier.
---
### Testing
- Manual API testing was performed using Postman
- Current version does not have automated tests