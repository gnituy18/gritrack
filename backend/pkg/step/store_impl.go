package step

import (
	"time"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"gritter/pkg/context"
	"gritter/pkg/document"
)

func New() Store {
	return &impl{
		doc: document.NewDocument(),
	}
}

var (
	uuidNewV4 = uuid.NewV4
)

type impl struct {
	doc document.Document
}

func (im *impl) Create(ctx context.Context, missionId string, s *Step) (string, error) {
	id := uuidNewV4().String()
	s.Id = id
	s.MissionId = missionId
	s.CreatedAt = time.Now().Unix()
	if err := im.doc.CreateOne(ctx, document.Step, s); err != nil {
		ctx.With(
			zap.Error(err),
			zap.Object("step", s),
		).Error("document.Document.CreateOne failed in step.Store.Create")
		return "", err
	}

	return id, nil
}
func (im *impl) Update(ctx context.Context, itemId string, s *Step) error {
	return nil
}
