package api

import (
	"net/http"

	routing "github.com/qiangxue/fasthttp-routing"

	"gritter/pkg/context"
	"gritter/pkg/user"
)

func MountUserRoutes(group *routing.RouteGroup, userStore user.Store) {
	handler := &userHandler{
		userStore: userStore,
	}

	group.Get("/current", handler.getCurrent)
}

type userHandler struct {
	userStore user.Store
}

func (uh *userHandler) getCurrent(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)

	val := rctx.Get("userId")
	userId, ok := val.(string)
	if !ok || userId == "" {
		JSON(rctx, http.StatusUnauthorized, nil)
		return nil
	}

	u, err := uh.userStore.Get(ctx, userId)
	if err == user.ErrNotFound {
		JSON(rctx, http.StatusNotFound, nil)
		return nil
	} else if err != nil {
		JSON(rctx, http.StatusInternalServerError, err.Error())
		return nil
	}

	JSON(rctx, http.StatusOK, u)
	return nil
}
