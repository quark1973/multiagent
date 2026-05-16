// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/quark1973/multiagent/utility/appconfig"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DSThinkChatModel     appconfig.ModelConfig  `json:"ds_think_chat_model,optional"`
	DSQuickChatModel     appconfig.ModelConfig  `json:"ds_quick_chat_model,optional"`
	DoubaoEmbeddingModel appconfig.ModelConfig  `json:"doubao_embedding_model,optional"`
	FileDir              string                 `json:"file_dir,optional"`
	MCPURL               string                 `json:"mcp_url,optional"`
	Feishu               appconfig.FeishuConfig `json:"feishu,optional"`
}

func (c Config) AppConfig() appconfig.Config {
	return appconfig.Config{
		DSThinkChatModel:     c.DSThinkChatModel,
		DSQuickChatModel:     c.DSQuickChatModel,
		DoubaoEmbeddingModel: c.DoubaoEmbeddingModel,
		FileDir:              c.FileDir,
		MCPURL:               c.MCPURL,
		Feishu:               c.Feishu,
	}
}
