// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"net/http"

	"github.com/quark1973/multiagent/internal/svc"
	"github.com/quark1973/multiagent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(r *http.Request) (resp *types.FileUploadRes, err error) {
	uploaded, err := l.svcCtx.KnowledgeService.UploadAndIndex(l.ctx, r)
	if err != nil {
		return nil, err
	}
	return &types.FileUploadRes{
		FileName: uploaded.FileName,
		FilePath: uploaded.FilePath,
		FileSize: uploaded.FileSize,
	}, nil
}
