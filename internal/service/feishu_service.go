package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/quark1973/multiagent/utility/appconfig"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type FeishuService struct {
	client       *lark.Client
	handler      http.HandlerFunc
	chatService  *ChatService
	aiOpsService *AIOpsService
	enabled      bool
}

func NewFeishuService(cfg appconfig.FeishuConfig, chatService *ChatService, aiOpsService *AIOpsService) *FeishuService {
	s := &FeishuService{
		chatService:  chatService,
		aiOpsService: aiOpsService,
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		return s
	}

	s.client = lark.NewClient(cfg.AppID, cfg.AppSecret)
	eventDispatcher := dispatcher.NewEventDispatcher(cfg.VerifyToken, cfg.EncryptKey).
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			chatID, prompt, ok := parseFeishuTextMessage(event)
			if !ok {
				return nil
			}

			go s.handleMessage(chatID, prompt)
			return nil
		})

	s.handler = httpserverext.NewEventHandlerFunc(eventDispatcher)
	s.enabled = true
	return s
}

func (s *FeishuService) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	if !s.enabled || s.handler == nil {
		http.Error(w, "feishu bot is not configured", http.StatusServiceUnavailable)
		return nil
	}
	s.handler(w, r)
	return nil
}

func (s *FeishuService) handleMessage(chatID string, prompt string) {
	ctx := context.Background()
	if err := s.sendText(ctx, chatID, "received, analyzing..."); err != nil {
		log.Printf("[feishu] send ack failed: %v", err)
	}

	var (
		result string
		err    error
	)
	if shouldRunAIOps(prompt) {
		result, _, err = s.aiOpsService.Analyze(ctx)
	} else {
		result, err = s.chatService.Chat(ctx, chatID, prompt)
	}
	if err != nil {
		result = fmt.Sprintf("agent failed: %v", err)
	}
	if result == "" {
		result = "agent returned empty response"
	}
	if err := s.sendText(ctx, chatID, result); err != nil {
		log.Printf("[feishu] send result failed: %v", err)
	}
}

func (s *FeishuService) sendText(ctx context.Context, chatID string, text string) error {
	contentBytes, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return err
	}
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(chatID).
			MsgType(larkim.MsgTypeText).
			Content(string(contentBytes)).
			Build()).
		Build()
	_, err = s.client.Im.Message.Create(ctx, req)
	return err
}

func parseFeishuTextMessage(event *larkim.P2MessageReceiveV1) (string, string, bool) {
	if event == nil || event.Event == nil || event.Event.Message == nil {
		return "", "", false
	}
	msg := event.Event.Message
	if msg.ChatId == nil || msg.Content == nil {
		return "", "", false
	}
	if msg.MessageType != nil && *msg.MessageType != "text" {
		return "", "", false
	}

	var body struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(*msg.Content), &body); err != nil {
		return "", "", false
	}
	prompt := strings.TrimSpace(body.Text)
	if prompt == "" {
		return "", "", false
	}
	return *msg.ChatId, prompt, true
}

func shouldRunAIOps(prompt string) bool {
	trimmed := strings.TrimSpace(prompt)
	return strings.HasPrefix(strings.ToLower(trimmed), "/aiops") || strings.Contains(trimmed, "告警分析")
}
