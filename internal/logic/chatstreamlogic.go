// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/quark1973/multiagent/internal/svc"
	"github.com/quark1973/multiagent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatStreamLogic {
	return &ChatStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatStreamLogic) ChatStream(req *types.ChatStreamReq, w http.ResponseWriter) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming is not supported")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	writeEvent := func(eventType string, data string) error {
		_, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, strings.ReplaceAll(data, "\n", "\\n"))
		if err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	if err := writeEvent("connected", `{"status":"connected"}`); err != nil {
		return err
	}
	if err := l.svcCtx.ChatService.Stream(l.ctx, req.Id, req.Question, func(chunk string) error {
		return writeEvent("message", chunk)
	}); err != nil {
		_ = writeEvent("error", err.Error())
		return err
	}
	return writeEvent("done", "Stream completed")
}
