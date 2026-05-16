// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"github.com/quark1973/multiagent/internal/config"
	"github.com/quark1973/multiagent/internal/service"
	"github.com/quark1973/multiagent/utility/appconfig"
	"github.com/quark1973/multiagent/utility/common"
)

type ServiceContext struct {
	Config           config.Config
	ChatService      *service.ChatService
	KnowledgeService *service.KnowledgeService
	AIOpsService     *service.AIOpsService
	FeishuService    *service.FeishuService
}

func NewServiceContext(c config.Config) *ServiceContext {
	appCfg := c.AppConfig()
	appconfig.Set(appCfg)
	if appCfg.FileDir != "" {
		common.FileDir = appCfg.FileDir
	}

	chatService := service.NewChatService()
	aiOpsService := service.NewAIOpsService()
	return &ServiceContext{
		Config:           c,
		ChatService:      chatService,
		KnowledgeService: service.NewKnowledgeService(),
		AIOpsService:     aiOpsService,
		FeishuService:    service.NewFeishuService(appCfg.Feishu, chatService, aiOpsService),
	}
}
