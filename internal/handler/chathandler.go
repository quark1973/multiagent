// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/quark1973/multiagent/internal/logic"
	"github.com/quark1973/multiagent/internal/svc"
	"github.com/quark1973/multiagent/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewChatLogic(r.Context(), svcCtx)
		resp, err := l.Chat(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
