package user

import (
	"errors"

	"gritter/pkg/auth"
	"gritter/pkg/context"
)

var (
	ErrNotFound = errors.New("user not found")
)

type Store interface {
	Create(ctx context.Context, u *User) (string, error)
	CreateByAuthResult(ctx context.Context, result *auth.Result) (string, error)
	Get(ctx context.Context, id string) (*User, error)
	GetByAuthResult(ctx context.Context, result *auth.Result) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id string) error
}
