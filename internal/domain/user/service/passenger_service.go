package service

import (
	"errors"
	"github.com/gavin/airport-pickup/internal/domain/user/entity"
)

type PassengerService struct{}

type CreatePassengerCmd struct {
	Name string
}

func (s *PassengerService) CreatePassenger(cmd *CreatePassengerCmd) (*entity.Passenger, error) {
	if cmd.Name == "" {
		return nil, errors.New("name required")
	}
	return &entity.Passenger{ID: "", Name: cmd.Name}, nil
}
