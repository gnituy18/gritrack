package auth

import (
	"errors"

	"gritter/pkg/context"
)

var (
	ErrTypeInvalid          = errors.New("type invalid")
	ErrTokenAudienceInvalid = errors.New("token audience invalid")
)

type Store interface {
	Auth(ctx context.Context, info *Info) (*Result, error)
}
