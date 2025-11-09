# Airport Pickup Order Matching System

## 1. Project Overview

The Airport Pickup Order Matching System is a backend service designed to efficiently match passengers and drivers for airport pickups. The system considers airport location, time slot, vehicle type, and pricing rules to optimize matches. The platform generates revenue from the margin between passenger bids and driver offers.

## 2. System Architecture

The system adopts Domain-Driven Design (DDD) principles, organizing code into distinct modules:
- **API Layer** (`api/`): Handles HTTP and gRPC requests.
- **Application Layer** (`internal/app/`): Coordinates business use cases.
- **Domain Layer** (`internal/domain/`): Contains core business logic and entities (order, settlement, user, eventbus).
- **Infrastructure Layer** (`pkg/`): Provides integrations (database, event bus, payments, Redis, HTTP utilities).

This separation ensures maintainability and scalability. The diagram below illustrates the architecture:

```
[API Layer] <-> [Application Layer] <-> [Domain Layer] <-> [Infrastructure Layer]
```

## 3. Tech Stack

- **Language:** Go 1.22
- **Frameworks:** Gin (HTTP), Gorm (ORM)
- **Database:** MySQL, Redis
- **Message Queue:** Kafka
- **Testing:** Go Test
- **Containerization:** Docker Compose

## 4. Run Instructions

```bash
# Clone repository
$ git clone xxx
$ cd airport-pickup

# Run with Docker Compose
$ docker compose up

# Or run locally
$ go run cmd/server/main.go
```
# db migrate
docker exec -i airport-mysql mysql -uairport -pairport airport < db/migrations/001_init_schema.sql

## 5. API Documentation

### 1. Create Passenger
- **POST** `/passengers`
- **Request Body:**
  ```json
  {
    "name": "Alice"
  }
  ```

### 2. Create Driver
- **POST** `/drivers`
- **Request Body:**
  ```json
  {
    "name": "Bob",
    "rating": 4.9
  }
  ```

### 3. Create Pickup Request
- **POST** `/pickup_requests`
- **Request Body:**
  ```json
  {
    "passenger_id": "174b032d1244ea6320a77041c034bd8f",
    "airport_code": "SFO",
    "vehicle_type": "sedan",
    "desired_time": "2025-11-05T10:00:00Z",
    "max_price_per_km": 2.5,
    "prefer_high_rating": true
  }
  ```

### 4. Create Driver Offer
- **POST** `/driver_offers`
- **Request Body:**
  ```json
  {
    "driver_id": "0bd803342d1661d5380c833f04929417",
    "airport_code": "SFO",
    "vehicle_type": "sedan",
    "available_from": "2025-11-05T09:00:00Z",
    "available_to": "2025-11-05T12:00:00Z",
    "price_per_km": 2.0
  }
  ```

### 5. List Bookings
- **GET** `/bookings`

### 6. Complete Booking
- **POST** `/bookings?id=ed6c04d6777b4d782f312519623fdf18`

## 6. Domain Model / Matching Logic

The matching algorithm works as follows:
1. Select drivers whose available time slots overlap with the passenger's requested time.
2. Filter drivers by vehicle type and (optionally) rating.
3. Choose the driver offering the lowest price that does not exceed the passenger's maximum bid.

**Pseudocode:**
```
for each passenger_request:
    candidates = find drivers where
        driver.available_time overlaps passenger.desired_time
        and driver.vehicle_type == passenger.vehicle_type
        and driver.price_per_km <= passenger.max_price_per_km
    if candidates:
        select driver with lowest price_per_km
        create match
```
