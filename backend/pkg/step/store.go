package step

import (
	"gritter/pkg/context"
)

type Store interface {
	Create(ctx context.Context, missionId string) (string error)
	Update(ctx context.Context, itemId string, item *Item) error
}
