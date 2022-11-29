package handler

import (
	"net/http"

	"github.com/cuczhangyi/coros-strava/internal/logic"
	"github.com/cuczhangyi/coros-strava/internal/svc"
	"github.com/cuczhangyi/coros-strava/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func StravaCallbackHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StravaCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewStravaCallbackLogic(r.Context(), ctx)
		resp, err := l.StravaCallback(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
