package api

import (
	"fmt"
	"net/http"

	routing "github.com/qiangxue/fasthttp-routing"
	"go.uber.org/zap"

	"gritter/pkg/context"
	"gritter/pkg/mission"
	"gritter/pkg/user"
)

func MountUserRoutes(group *routing.RouteGroup, userStore user.Store, missionStore mission.Store) {
	handler := &userHandler{
		userStore: userStore,
		missionStore: missionStore,
	}

	group.Get("/current", handler.getCurrent)
	group.Get("/<userId>/missionName/<missionName>", handler.getMission)
}

type userHandler struct {
	userStore user.Store
	missionStore mission.Store
}

func (uh *userHandler) getMission(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)
	userId := rctx.Param("userId")
	missionName := rctx.Param("missionName")

	fmt.Println(userId, missionName)

	m, err := uh.missionStore.GetByUserMissionName(ctx, userId, missionName)
	if err == mission.ErrNotFound {
		JSON(rctx, http.StatusNotFound, err.Error())
		return nil
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("userId", userId),
			zap.String("missionName", missionName),
		).Error("userHandler.missionStore.GetByUserMissionName failed in userHandler.getMission")
		JSON(rctx, http.StatusInternalServerError, nil)
		return nil
	}

	repr := missionToRepr(m)

	JSON(rctx, http.StatusOK, repr)
	return nil
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

	userRepr := userToRepr(u)

	JSON(rctx, http.StatusOK, userRepr)
	return nil
}
