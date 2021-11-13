package api

import (
	"encoding/json"
	"net/http"

	routing "github.com/qiangxue/fasthttp-routing"
	"go.uber.org/zap"

	"gritter/pkg/auth"
	"gritter/pkg/context"
	"gritter/pkg/user"
)

func MountAuthRoutes(group *routing.RouteGroup, authStore auth.Store, userStore user.Store) {
	handler := &authHandler{
		authStore: authStore,
		userStore: userStore,
	}

	group.Post("", handler.auth)
}

type authHandler struct {
	authStore auth.Store
	userStore user.Store
}

func (ah *authHandler) auth(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)

	body := rctx.Request.Body()
	info := &auth.Info{}
	if err := json.Unmarshal(body, info); err != nil {
		JSON(rctx, http.StatusBadRequest, err.Error())
		return nil
	}

	result, err := ah.authStore.Auth(ctx, info)
	if err != nil {
		ctx.With(zap.Error(err)).Error("authHandler.authStore.Auth failed in authHandler.auth")
		JSON(rctx, http.StatusInternalServerError, nil)
		return err
	}

	// search for existed user
	u, err := ah.userStore.GetByAuthResult(ctx, result)
	if err != nil && err != user.ErrNotFound {
		ctx.With(zap.Error(err)).Error("authHandler.userStore.GetByAuthResult failed in authHandler.auth")
		JSON(rctx, http.StatusInternalServerError, nil)
		return err
	}

	// create user if not exist
	var userId string
	if err == user.ErrNotFound {
		userId, err = ah.userStore.CreateByAuthResult(ctx, result)
		if err != nil {
			ctx.With(zap.Error(err)).Error("authHandler.userStore.CreateByAuthResult failed in authHandler.auth")
			JSON(rctx, http.StatusInternalServerError, nil)
			return err
		}
	} else {
		userId = u.Id
	}

	// store userId in session
	store, err := sess.Get(rctx.RequestCtx)
	if err != nil {
		ctx.Error("session.Session.Get failed in authHandler.auth")
		return err
	}

	store.Set("userId", userId)
	saveStore(rctx, store)

	Redirect(rctx, "http://localhost:3000")

	return nil
}
