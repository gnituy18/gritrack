package auth

import (
	"errors"

	"gritter/pkg/context"
)

var (
	ErrTypeInvalid = errors.New("type invalid")
)

type Store interface {
	Auth(ctx context.Context, info *Info) (*Result, error)
}
