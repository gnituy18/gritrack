package step

import (
	"time"

	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
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

func (im *impl) Create(ctx context.Context, s *Step) (string, error) {
	id := uuidNewV4().String()
	s.Id = id
	s.CreatedAt = time.Now().Unix()
	if s.Items == nil {
		s.Items = []*Item{}
	}

	if err := im.doc.CreateOne(ctx, document.Step, s); err != nil {
		ctx.With(
			zap.Error(err),
			zap.Object("step", s),
		).Error("document.Document.CreateOne failed in step.Store.Create")
		return "", err
	}

	return id, nil
}
func (im *impl) Update(ctx context.Context, s *Step) error {
	q := bson.M{
		"id":        s.Id,
		"missionId": s.MissionId,
	}
	updater := bson.M{
		"summary": s.Summary,
		"items":   s.Items,
	}
	if err := im.doc.UpdateOne(ctx, document.Step, q, updater); err == document.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("id", s.Id),
		).Error("document.Document.UpdateOne failed in step.Store.Update")
		return err
	}

	return nil
}
