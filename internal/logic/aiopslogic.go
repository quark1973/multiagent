// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"github.com/quark1973/multiagent/internal/svc"
	"github.com/quark1973/multiagent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AIOpsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAIOpsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AIOpsLogic {
	return &AIOpsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AIOpsLogic) AIOps() (resp *types.AIOpsRes, err error) {
	result, detail, err := l.svcCtx.AIOpsService.Analyze(l.ctx)
	if err != nil {
		return nil, err
	}
	return &types.AIOpsRes{
		Result: result,
		Detail: detail,
	}, nil
}
