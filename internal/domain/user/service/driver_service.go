package service

import (
	"errors"
	"github.com/gavin/airport-pickup/internal/domain/user/entity"
)

type DriverService struct{}

type CreateDriverCmd struct {
	Name   string
	Rating float64
}

func (s *DriverService) CreateDriver(cmd *CreateDriverCmd) (*entity.Driver, error) {
	if cmd.Name == "" {
		return nil, errors.New("name required")
	}
	if cmd.Rating < 0 || cmd.Rating > 5 {
		return nil, errors.New("invalid rating")
	}
	return &entity.Driver{ID: "", Name: cmd.Name, Rating: cmd.Rating}, nil
}
