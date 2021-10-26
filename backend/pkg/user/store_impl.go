package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"gritter/pkg/auth"
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

func (im *impl) Create(ctx context.Context, u *User) (string, error) {
	id := uuidNewV4().String()
	u.ID = id
	if err := im.doc.CreateOne(ctx, document.User, u); err != nil {
		ctx.With(
			zap.Error(err),
			zap.Object("user", u),
		).Error("document.Document.CreateOne failed in user.Store.Create")
		return "", err
	}

	return id, nil
}

func (im *impl) CreateByAuthResult(ctx context.Context, result *auth.Result) (string, error) {
	var u *User
	switch result.Type {
	case auth.TypeGoogle:
		u = im.authGoogleToUser(ctx, result.Google)
	default:
		return "", errors.New("auth type invalid")
	}

	id := uuidNewV4().String()
	u.ID = id

	if err := im.doc.CreateOne(ctx, document.User, u); err != nil {
		ctx.With(
			zap.Error(err),
			zap.Object("user", u),
		).Error("document.Document.CreateOne failed in user.Store.Create")
		return "", err
	}

	return id, nil
}

func (im *impl) authGoogleToUser(ctx context.Context, result *auth.ResultGoogle) *User {
	return &User{
		Email:        result.TokenInfo.Email,
		GoogleUserID: result.TokenInfo.UserId,
	}
}

func (im *impl) Get(ctx context.Context, id string) (*User, error) {
	q := bson.M{
		"id":      id,
		"deleted": false,
	}
	u := &User{}
	if err := im.doc.GetOne(ctx, document.User, q, u); err == document.ErrNotFound {
		return nil, ErrNotFound
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("id", id),
		).Error("document.Document.GetOne failed in user.Store.Get")
		return nil, err
	}

	return u, nil
}

func (im *impl) GetByAuthResult(ctx context.Context, result *auth.Result) (*User, error) {
	var q bson.M
	switch result.Type {
	case auth.TypeGoogle:
		q = im.getQueryByGoogle(ctx, result.Google)
	default:
		return nil, errors.New("invalid auth type")
	}

	u := &User{}
	if err := im.doc.GetOne(ctx, document.User, q, u); err == document.ErrNotFound {
		return nil, ErrNotFound
	} else if err != nil {
		ctx.With(
			zap.Error(err),
		).Error("document.Document.GetOne failed in user.Store.GetByAuthResult")
		return nil, err
	}

	fmt.Println(u)

	return u, nil
}

func (im *impl) getQueryByGoogle(ctx context.Context, result *auth.ResultGoogle) bson.M {
	return bson.M{
		"googleUserID": result.TokenInfo.UserId,
	}
}

func (im *impl) Update(ctx context.Context, u *User) error {
	updater := bson.M{
		"name":  u.Name,
		"intro": u.Intro,
	}
	if err := im.doc.UpdateOne(ctx, document.User, bson.M{"id": u.ID}, updater); err == document.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("id", u.ID),
		).Error("document.Document.UpdateOne failed in user.Store.Update")
		return err
	}

	return nil
}

func (im *impl) Delete(ctx context.Context, id string) error {
	updater := bson.M{
		"deleted": true,
	}
	if err := im.doc.UpdateOne(ctx, document.User, bson.M{"id": id}, updater); err == document.ErrNotFound {
		return nil
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("id", id),
		).Error("document.Document.UpdateOne failed in user.Store.Update")
		return err
	}

	return nil
}
