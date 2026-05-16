// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"github.com/quark1973/multiagent/internal/svc"
	"github.com/quark1973/multiagent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatLogic) Chat(req *types.ChatReq) (resp *types.ChatRes, err error) {
	answer, err := l.svcCtx.ChatService.Chat(l.ctx, req.Id, req.Question)
	if err != nil {
		return nil, err
	}
	return &types.ChatRes{Answer: answer}, nil
}
