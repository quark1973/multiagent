// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"net/http"

	"github.com/quark1973/multiagent/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FeishuEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeishuEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeishuEventLogic {
	return &FeishuEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeishuEventLogic) FeishuEvent(w http.ResponseWriter, r *http.Request) error {
	return l.svcCtx.FeishuService.ServeHTTP(w, r)
}
