package api

import (
	"encoding/json"
	"net/http"

	"github.com/qiangxue/fasthttp-routing"
	"go.uber.org/zap"

	"gritter/pkg/context"
	"gritter/pkg/mission"
)

func MountMissionRoutes(group *routing.RouteGroup, missionStore mission.Store) {
	handler := &missionHandler{
		missionStore: missionStore,
	}

	group.Get("", handler.getMissionsByUserId)
	group.Post("", handler.createMission)
	group.Get("/<missionId>", handler.getMission)

	mustOwnMission := createMustOwnMission(missionStore)
	group.Use(mustOwnMission)
	group.Put("/<missionId>", handler.updateMission)
	group.Delete("/<missionId>", handler.deleteMission)
}

type missionHandler struct {
	missionStore mission.Store
}

type createMissionBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (mh *missionHandler) getMissionsByUserId(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)
	userId, ok := rctx.Get("userId").(string)
	if !ok || userId == "" {
		ctx.Error("rctx.Get userId failed in missionHandler.getMissionsByUserId")
		JSON(rctx, http.StatusInternalServerError, nil)
		return nil
	}
	missions, err := mh.missionStore.GetByUser(ctx, userId)
	if err != nil {
		ctx.Error("missionHandler.missionStore.GetByUser failed in missionHandler.getMissionsByUserId")
		JSON(rctx, http.StatusInternalServerError, nil)
		return nil
	}

	missionsRepr := []*missionRepr{}
	for _, m := range missions {
		missionsRepr = append(missionsRepr, missionToRepr(m))
	}

	JSON(rctx, http.StatusOK, missionsRepr)
	return nil
}

func (mh *missionHandler) createMission(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)
	userId, ok := rctx.Get("userId").(string)
	if !ok || userId == "" {
		ctx.Error("rctx.Get userId failed in missionHandler.createMission")
		JSON(rctx, http.StatusInternalServerError, nil)
		return nil
	}

	body := &createMissionBody{}
	if err := json.Unmarshal(rctx.Request.Body(), body); err != nil {
		JSON(rctx, http.StatusBadRequest, err.Error())
		return nil
	}

	m := &mission.Mission{
		UserId:      userId,
		Name:        body.Name,
		Description: body.Description,
	}
	id, err := mh.missionStore.Create(ctx, m)
	if err != nil {
		ctx.With(zap.Error(err)).Error("missionHandler.missionStore.Create failed in missionHandler.createMission")
		JSON(rctx, http.StatusInternalServerError, nil)
		return err
	}

	JSON(rctx, http.StatusCreated, id)
	return nil
}

func (mh *missionHandler) getMission(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)
	id := rctx.Param("missionId")

	m, err := mh.missionStore.Get(ctx, id)
	if err == mission.ErrNotFound {
		JSON(rctx, http.StatusNotFound, err.Error())
		return nil
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("missionId", id),
		).Error("missionHandler.missionStore.Get failed in missionHandler.getMission")
		JSON(rctx, http.StatusInternalServerError, nil)
		return err
	}

	repr := missionToRepr(m)

	JSON(rctx, http.StatusOK, repr)
	return nil
}

func (mh *missionHandler) updateMission(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)
	id := rctx.Param("missionId")

	body := rctx.Request.Body()
	m := &mission.Mission{}
	if err := json.Unmarshal(body, m); err != nil {
		JSON(rctx, http.StatusBadRequest, err.Error())
		return nil
	}

	m.Id = id

	err := mh.missionStore.Update(ctx, m)
	if err == mission.ErrNotFound {
		JSON(rctx, http.StatusNotFound, err.Error())
		return nil
	} else if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("missionId", id),
		).Error("missionHandler.missionStore.Update failed in missionHandler.updateMission")
		JSON(rctx, http.StatusInternalServerError, nil)
		return err
	}

	JSON(rctx, http.StatusNoContent, nil)
	return nil
}

func (mh *missionHandler) deleteMission(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)
	id := rctx.Param("missionId")

	err := mh.missionStore.Delete(ctx, id)
	if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("missionId", id),
		).Error("missionHandler.missionStore.Delete failed in missionHandler.deleteMission")
		JSON(rctx, http.StatusInternalServerError, nil)
		return err
	}

	JSON(rctx, http.StatusNoContent, nil)
	return nil
}
