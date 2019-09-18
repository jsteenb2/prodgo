package users

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Store interface {
	Create(ctx context.Context, u User) (User, error)
	Read(ctx context.Context, id string) (User, error)
}

type Service struct {
	store Store
}

func (s *Service) Create(ctx context.Context, u User) (User, error) {
	const op = "UserService.CreateUser"

	newUser, err := s.store.Create(ctx, u)
	if err != nil {
		return User{}, errors.Wrap(err, op)
	}

	return newUser, nil
}

func (s *Service) Read(ctx context.Context, id string) (User, error) {
	const op = "UserService.Read"

	user, err := s.store.Read(ctx, id)
	if err != nil {
		return User{}, errors.Wrap(err, op)
	}

	return user, nil
}
