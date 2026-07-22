package domain

import "github.com/google/uuid"

type Repository interface {
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	FindAll() ([]User, error)
	FindAllWithRole() ([]User, error)
}

type RoleRepository interface {
	FindRoleByID(id uuid.UUID) (*Role, error)
}

type Role struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
