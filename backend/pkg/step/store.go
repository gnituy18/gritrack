package step

import (
	"gritter/pkg/context"
)

type Store interface {
	Create(ctx context.Context, missionId string, step *Step) (string, error)
	Update(ctx context.Context, itemId string, step *Step) error
}
