package step

import (
	"errors"

	"gritter/pkg/context"
)

var (
	ErrNotFound = errors.New("step not found")
)

type Store interface {
	IsStepTimeExists(ctx context.Context, ts int64) (bool, error)
	GetByMissionId(ctx context.Context, missionId string, offset, limit int64) ([]*Step, error)
	Create(ctx context.Context, step *Step) (string, error)
	Update(ctx context.Context, step *Step) error
}
