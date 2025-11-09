package entity

import "time"

type Driver struct {
	ID        string
	Name      string
	Rating    float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
