package models

import (
	"context"
	"github.com/quark1973/multiagent/utility/appconfig"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

func OpenAIForDeepSeekV31Think(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	cfg, err := appconfig.Get()
	if err != nil {
		return nil, err
	}
	modelCfg := cfg.DSThinkChatModel
	config := &openai.ChatModelConfig{
		Model:   modelCfg.Model,
		APIKey:  modelCfg.APIKey,
		BaseURL: modelCfg.BaseURL,
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func OpenAIForDeepSeekV3Quick(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	cfg, err := appconfig.Get()
	if err != nil {
		return nil, err
	}
	modelCfg := cfg.DSQuickChatModel
	config := &openai.ChatModelConfig{
		Model:   modelCfg.Model,
		APIKey:  modelCfg.APIKey,
		BaseURL: modelCfg.BaseURL,
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
