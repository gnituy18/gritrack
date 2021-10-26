package api

import (
	"encoding/json"
	"net/http"

	"github.com/qiangxue/fasthttp-routing"

	"gritter/pkg/log"
)

var (
	contentTypeJSON = "application/json"
)

func JSON(rctx *routing.Context, status int, body interface{}) {
	rctx.SetContentType(contentTypeJSON)
	rctx.SetStatusCode(status)

	bodyBytes := []byte{}
	if body != nil {
		bytes, err := json.Marshal(body)
		if err != nil {
			log.Global().Error("json.Marshal failed in api.JSON")
			rctx.Error("json.Marshal failed in api.JSON", http.StatusInternalServerError)
			return
		}

		bodyBytes = bytes
	}

	rctx.SetBody(bodyBytes)
}
