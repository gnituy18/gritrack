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

func (im *impl) GetByMissionId(ctx context.Context, missionId string, offset, limit int64) ([]*Step, bool, error) {
	q := bson.M{
		"missionId": missionId,
	}
	sort := bson.D{bson.E{Key: "time", Value: -1}}
	steps := []*Step{}
	err := im.doc.Search(ctx, document.Step, offset, limit+1, q, sort, &steps)
	if err != nil {
		ctx.With(zap.Error(err)).Error("document.Document.Search failed in step.Store.GetByMissionId")
		return nil, false, err
	}

	if len(steps) == 0 {
		return steps, false, nil
	}

	if int64(len(steps)) > limit {
		return steps[:len(steps)-1], true, nil
	}

	return steps, false, nil
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

func (im *impl) IsStepTimeExists(ctx context.Context, missionId string, ts int64) (bool, error) {
	q := bson.M{
		"missionId": missionId,
		"time":      ts,
	}

	if err := im.doc.GetOne(ctx, document.Step, q, &Step{}); err == document.ErrNotFound {
		return false, nil
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.Int64("time", ts),
		).Error("document.Document.GetOne failed in step.Store.IsStepDateExists")
		return false, err
	}

	return true, nil
}
