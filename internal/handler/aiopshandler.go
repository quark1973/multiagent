// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/quark1973/multiagent/internal/logic"
	"github.com/quark1973/multiagent/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AIOpsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAIOpsLogic(r.Context(), svcCtx)
		resp, err := l.AIOps()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
