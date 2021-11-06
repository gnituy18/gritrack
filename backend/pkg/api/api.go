package api

import (
	"net/http"

	routing "github.com/qiangxue/fasthttp-routing"

	"gritter/pkg/auth"
	"gritter/pkg/log"
	"gritter/pkg/mission"
	"gritter/pkg/user"
)

func Router() *routing.Router {
	authStore := auth.New()
	missionStore := mission.New()
	userStore := user.New()

	root := routing.New()
	root.Use(corsHeader)
	root.Use(injectContext)
	root.Use(logRequest)

	// public routes
	public := root.Group("/api/v1")

	public.Get("/health", func(rctx *routing.Context) error {
		rctx.WriteString("healthy")
		log.Global().Info("healthy")
		return nil
	})

	authGroup := public.Group("/auth")
	MountAuthRoutes(authGroup, authStore, userStore)

	// authenticated routes
	authenticated := root.Group("")
	authenticated.Use(mustAuthUser)

	missionGroup := authenticated.Group("/mission")
	MountMissionRoutes(missionGroup, missionStore)

	root.Any("*", func(rctx *routing.Context) error {
		JSON(rctx, http.StatusNotFound, nil)
		return nil
	})

	return root
}
