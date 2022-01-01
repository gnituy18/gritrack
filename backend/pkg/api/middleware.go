package api

import (
	"net/http"
	"os"
	"time"

	"github.com/fasthttp/session/v2"
	"github.com/fasthttp/session/v2/providers/redis"
	routing "github.com/qiangxue/fasthttp-routing"
	"go.uber.org/zap"

	"gritter/pkg/context"
	"gritter/pkg/log"
	"gritter/pkg/mission"
)

var (
	sess *session.Session
)

func init() {
	cfg := session.NewDefaultConfig()
	sess = session.New(cfg)

	provider, err := redis.New(redis.Config{
		KeyPrefix:   "session",
		Addr:        os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		PoolSize:    10,
		IdleTimeout: 24 * time.Hour,
	})
	if err != nil {
		log.Global().Fatal(err)
	}

	if err := sess.SetProvider(provider); err != nil {
		log.Global().Fatal(err)
	}
}

func injectContext(rctx *routing.Context) error {
	rctx.Set("ctx", context.Background())
	rctx.Next()
	return nil
}

func logRequest(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)
	ctx.With(zap.String("req", rctx.String())).Info("req")
	rctx.Next()
	return nil
}

func corsHeader(rctx *routing.Context) error {
	rctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	rctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	if rctx.Request.Header.IsOptions() {
		rctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE")
		rctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type")
		rctx.Response.Header.Set("Access-Control-Max-Age", "86400")
		rctx.SetContentType("text/plain")
		rctx.SetStatusCode(http.StatusOK)
		rctx.Abort()
		return nil
	}

	rctx.Next()
	return nil
}

func mustAuthUser(rctx *routing.Context) error {
	ctx := rctx.Get("ctx").(context.Context)

	store, err := sess.Get(rctx.RequestCtx)
	if err != nil {
		ctx.Error("session.Session.Get failed in api.GetAuthUser")
		return err
	}
	val := store.Get("userId")
	userId, ok := val.(string)
	if !ok || userId == "" {
		JSON(rctx, http.StatusUnauthorized, nil)
		rctx.Abort()
		return nil
	}

	rctx.Set("userId", userId)

	rctx.Next()

	return nil
}

func createMustOwnMission(ms mission.Store) func(rctx *routing.Context) error {
	return func(rctx *routing.Context) error {
		ctx := rctx.Get("ctx").(context.Context)
		missionId := rctx.Param("missionId")
		userId := rctx.Get("userId").(string)
		if owned, err := ms.OwnedBy(ctx, missionId, userId); err != nil {
			ctx.Error("missionStore.OwnedBy failed in api.createMustHasMission")
			return err
		} else if !owned {
			JSON(rctx, http.StatusForbidden, nil)
			rctx.Abort()
			return nil
		}

		rctx.Next()
		return nil
	}
}

func saveStore(rctx *routing.Context, store *session.Store) {
	if err := sess.Save(rctx.RequestCtx, store); err != nil {
		log.Global().Error(err)
	}
}
