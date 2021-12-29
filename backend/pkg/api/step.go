package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gritter/pkg/context"
	"gritter/pkg/mission"
	"gritter/pkg/step"

	routing "github.com/qiangxue/fasthttp-routing"
)

func MountStepRoutes(group *routing.RouteGroup, stepStore step.Store, missionStore mission.Store) {
	handler := &stepHandler{
		stepStore:    stepStore,
		missionStore: missionStore,
	}

	group.Get("", handler.getByMissionId)
	group.Post("", handler.createStep)
	group.Put("/<stepId>", handler.updateStep)
}

type stepHandler struct {
	stepStore    step.Store
	missionStore mission.Store
}

type stepBody struct {
	Summary string     `json:"summary"`
	Items   step.Items `json:"items"`
}

func (sh *stepHandler) getByMissionId(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)
	missionId := rctx.Param("missionId")
	offset, err := strconv.ParseInt(string(rctx.QueryArgs().Peek("offset")), 10, 64)
	if err != nil {
		JSON(rctx, http.StatusBadRequest, nil)
	}
	limit, err := strconv.ParseInt(string(rctx.QueryArgs().Peek("limit")), 10, 64)
	if err != nil {
		JSON(rctx, http.StatusBadRequest, nil)
	}

	steps, err := sh.stepStore.GetByMissionId(ctx, missionId, offset, limit)
	reprs := []*stepRepr{}
	for _, s := range steps {
		reprs = append(reprs, stepToRepr(s))
	}

	JSON(rctx, http.StatusOK, reprs)
	return nil
}

func (sh *stepHandler) createStep(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)
	userId, ok := rctx.Get("userId").(string)
	if !ok || userId == "" {
		ctx.Error("rctx.Get userId failed in stepHandler.createStep")
		JSON(rctx, http.StatusInternalServerError, nil)
		return nil
	}

	missionId := rctx.Param("missionId")
	owned, err := sh.missionStore.OwnedBy(ctx, missionId, userId)
	if err != nil {
		ctx.Error("stepHandler.missionStore.OwnedBy failed in stepHandler.createStep")
		JSON(rctx, http.StatusInternalServerError, nil)
		return nil
	} else if !owned {
		JSON(rctx, http.StatusForbidden, nil)
		return nil
	}

	body := &stepBody{}
	if err := json.Unmarshal(rctx.Request.Body(), body); err != nil {
		JSON(rctx, http.StatusBadRequest, err.Error())
		return nil
	}

	step := &step.Step{
		MissionId: missionId,
		Summary:   body.Summary,
		Items:     body.Items,
	}
	stepId, err := sh.stepStore.Create(ctx, step)
	if err != nil {
		ctx.Error("stepHandler.stepStore.Create failed in stepHandler.createStep")
		JSON(rctx, http.StatusInternalServerError, nil)
		return nil
	}

	JSON(rctx, http.StatusCreated, stepId)
	return nil
}

func (sh *stepHandler) updateStep(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)
	userId, ok := rctx.Get("userId").(string)
	if !ok || userId == "" {
		ctx.Error("rctx.Get userId failed in stepHandler.updateStep")
		JSON(rctx, http.StatusInternalServerError, nil)
		return nil
	}

	missionId := rctx.Param("missionId")
	owned, err := sh.missionStore.OwnedBy(ctx, missionId, userId)
	if err != nil {
		ctx.Error("stepHandler.missionStore.OwnedBy failed in stepHandler.updateStep")
		JSON(rctx, http.StatusInternalServerError, nil)
		return nil
	} else if !owned {
		JSON(rctx, http.StatusForbidden, nil)
		return nil
	}

	body := &stepBody{}
	if err := json.Unmarshal(rctx.Request.Body(), body); err != nil {
		JSON(rctx, http.StatusBadRequest, err.Error())
		return nil
	}

	stepId := rctx.Param("stepId")
	s := &step.Step{
		Id:        stepId,
		MissionId: missionId,
		Summary:   body.Summary,
		Items:     body.Items,
	}
	err = sh.stepStore.Update(ctx, s)
	if err == step.ErrNotFound {
		JSON(rctx, http.StatusNotFound, nil)
		return nil
	} else if err != nil {
		ctx.Error("stepHandler.stepStore.Update failed in stepHandler.updateStep")
		JSON(rctx, http.StatusInternalServerError, nil)
		return nil
	}

	JSON(rctx, http.StatusOK, nil)
	return nil
}
