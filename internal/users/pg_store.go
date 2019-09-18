package users

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/alexcesaro/statsd"
	"github.com/jmoiron/sqlx"
)

type PGStore struct {
	db      *sqlx.DB
	metrics *statsd.Client
}

func NewPGStore(db *sqlx.DB) *PGStore {
	return &PGStore{
		db: db,
	}
}

type userEntity struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s *PGStore) Create(ctx context.Context, user User) (User, error) {
	const insertSQL = `
		INSERT INTO users (id, email, first_name, last_name, password)
		VALUES (:id, :email, :first_name, :last_name, :password)
		RETURNING *
	`

	query, args, err := s.db.BindNamed(insertSQL, modelToEnt(user))
	if err != nil {
		return User{}, err
	}

	var u userEntity
	if err := s.db.GetContext(ctx, &u, query, args...); err != nil {
		return User{}, err
	}

	return entToModel(u), nil
}

type Wherer func() (string, string)

func ByID(id string) Wherer {
	return func() (string, string) {
		return "id = ?", id
	}
}

func ByEmail(email string) Wherer {
	return func() (string, string) {
		return "email = ?", email
	}
}

func (s *PGStore) Read(ctx context.Context, wd Wherer) (User, error) {
	if wd == nil {
		return User{}, errors.New("invalid Wherer submitted")
	}

	clause, arg := wd()
	readSQL := fmt.Sprintf("SELECT * FROM users WHERE %s", clause)
	readSQL = s.db.Rebind(readSQL)

	var userEnt userEntity
	if err := s.db.GetContext(ctx, &userEnt, readSQL, arg); err != nil {
		return User{}, err
	}

	return entToModel(userEnt), nil
}

func (s *PGStore) Update(ctx context.Context, user User) (User, error) {
	const updateSQL = `
		UPDATE users SET email = :email, password = :password, first_name = :first_name, last_name = :last_name, updated_at = NOW()
		WHERE id = :id
		RETURNING *
	`

	query, args, err := s.db.BindNamed(updateSQL, modelToEnt(user))
	if err != nil {
		return User{}, err
	}

	var u userEntity
	if err := s.db.GetContext(ctx, &u, query, args...); err != nil {
		return User{}, err
	}

	return entToModel(u), nil
}

func (s *PGStore) Delete(ctx context.Context, wd Wherer) (string, error) {
	if wd == nil {
		return "", errors.New("invalid Wherer submitted")
	}

	clause, arg := wd()
	readSQL := fmt.Sprintf("DELETE FROM users WHERE %s", clause)
	readSQL = s.db.Rebind(readSQL)

	var userEnt userEntity
	if err := s.db.GetContext(ctx, &userEnt, readSQL, arg); err != nil {
		return "", err
	}

	return userEnt.ID, nil
}

func modelToEnt(u User) userEntity {
	return userEntity{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func entToModel(u userEntity) User {
	return User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
