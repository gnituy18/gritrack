package step

import (
	"errors"

	"gritter/pkg/context"
)

var (
	ErrNotFound = errors.New("step not found")
)

type Store interface {
	Create(ctx context.Context, step *Step) (string, error)
	Update(ctx context.Context, step *Step) error
}
