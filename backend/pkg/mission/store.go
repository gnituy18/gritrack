package mission

import (
	"errors"

	"gritter/pkg/context"
)

var (
	ErrNotFound = errors.New("mission not found")
)

type Store interface {
	OwnedBy(ctx context.Context, missonId, userId string) (bool, error)
	GetByUser(ctx context.Context, userId string) ([]*Mission, error)
	GetByUserMissionName(ctx context.Context, userId, missionName string) (*Mission, error)
	Create(ctx context.Context, m *Mission) (string, error)
	Get(ctx context.Context, id string) (*Mission, error)
	Update(ctx context.Context, m *Mission) error
	Delete(ctx context.Context, id string) error
}
