package mission

import (
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"gritter/pkg/config"
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

func (im *impl) OwnedBy(ctx context.Context, missionId, userId string) (bool, error) {
	q := bson.M{
		"id":      missionId,
		"userId":  userId,
		"deleted": false,
	}
	m := &Mission{}
	if err := im.doc.GetOne(ctx, document.Mission, q, m); err == document.ErrNotFound {
		return false, ErrNotFound
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("id", missionId),
			zap.String("userId", userId),
		).Error("document.Document.GetOne failed in mission.Store.OwnedBy")
		return false, err
	}

	return true, nil
}

func (im *impl) GetByUser(ctx context.Context, userId string) ([]*Mission, error) {
	q := bson.M{
		"userId": userId,
	}
	// Add sort by create time
	missions := []*Mission{}
	if err := im.doc.Search(ctx, document.Mission, 0, config.MissionMaxCount, q, nil, &missions); err != nil {
		ctx.With(zap.Error(err)).Error("document.Document.Search failed in mission.Store.GetByUser")
		return nil, err
	}
	return missions, nil
}

func (im *impl) Create(ctx context.Context, m *Mission) (string, error) {
	id := uuidNewV4().String()
	m.Id = id
	if err := im.doc.CreateOne(ctx, document.Mission, m); err != nil {
		ctx.With(
			zap.Error(err),
			zap.Object("mission", m),
		).Error("document.Document.CreateOne failed in mission.Store.Create")
		return "", err
	}

	return id, nil
}

func (im *impl) Get(ctx context.Context, id string) (*Mission, error) {
	q := bson.M{
		"id":      id,
		"deleted": false,
	}
	m := &Mission{}
	if err := im.doc.GetOne(ctx, document.Mission, q, m); err == document.ErrNotFound {
		return nil, ErrNotFound
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("id", id),
		).Error("document.Document.GetOne failed in mission.Store.Get")
		return nil, err
	}

	return m, nil
}

func (im *impl) Update(ctx context.Context, m *Mission) error {
	updater := bson.M{
		"name":        m.Name,
		"description": m.Description,
	}
	if err := im.doc.UpdateOne(ctx, document.Mission, bson.M{"id": m.Id}, updater); err == document.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("id", m.Id),
		).Error("document.Document.UpdateOne failed in mission.Store.Update")
		return err
	}

	return nil
}

func (im *impl) Delete(ctx context.Context, id string) error {
	updater := bson.M{
		"deleted": true,
	}
	if err := im.doc.UpdateOne(ctx, document.Mission, bson.M{"id": id}, updater); err == document.ErrNotFound {
		return nil
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("id", id),
		).Error("document.Document.UpdateOne failed in mission.Store.Update")
		return err
	}

	return nil
}
