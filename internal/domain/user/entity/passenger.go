package entity

import "time"

type Passenger struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
