package domain

import "github.com/google/uuid"

type Repository interface {
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
}
