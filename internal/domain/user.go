package domain

import (
	"context"
	"time"
)

type User struct {
	ID          string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
}

type UserList map[string]User

type UserUsecase interface {
	Create(ctx context.Context, displayName, email string) (User, error)
	GetByID(ctx context.Context, userID string) (User, error)
	GetAll(ctx context.Context) (UserList, error)
	Update(ctx context.Context, userID, displayName string) error
	Delete(ctx context.Context, userID string) error
}

type UserRepository interface {
	Create(ctx context.Context, displayName, email string) (User, error)
	GetByID(ctx context.Context, userID string) (User, error)
	GetAll(ctx context.Context) (UserList, error)
	Update(ctx context.Context, userID, displayName string) error
	Delete(ctx context.Context, userID string) error
}
