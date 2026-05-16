package appconfig

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type ModelConfig struct {
	APIKey  string `json:"api_key" yaml:"api_key"`
	BaseURL string `json:"base_url" yaml:"base_url"`
	Model   string `json:"model" yaml:"model"`
}

type FeishuConfig struct {
	AppID       string `json:"app_id" yaml:"app_id"`
	AppSecret   string `json:"app_secret" yaml:"app_secret"`
	VerifyToken string `json:"verify_token" yaml:"verify_token"`
	EncryptKey  string `json:"encrypt_key" yaml:"encrypt_key"`
}

type Config struct {
	DSThinkChatModel     ModelConfig  `json:"ds_think_chat_model" yaml:"ds_think_chat_model"`
	DSQuickChatModel     ModelConfig  `json:"ds_quick_chat_model" yaml:"ds_quick_chat_model"`
	DoubaoEmbeddingModel ModelConfig  `json:"doubao_embedding_model" yaml:"doubao_embedding_model"`
	FileDir              string       `json:"file_dir" yaml:"file_dir"`
	MCPURL               string       `json:"mcp_url" yaml:"mcp_url"`
	Feishu               FeishuConfig `json:"feishu" yaml:"feishu"`
}

var (
	mu     sync.RWMutex
	config Config
	loaded bool
)

func Set(c Config) {
	mu.Lock()
	defer mu.Unlock()
	config = c
	loaded = true
}

func Get() (Config, error) {
	mu.RLock()
	if loaded {
		defer mu.RUnlock()
		return config, nil
	}
	mu.RUnlock()

	if err := Load(DefaultPath()); err != nil {
		return Config{}, err
	}

	mu.RLock()
	defer mu.RUnlock()
	return config, nil
}

func Load(path string) error {
	if path == "" {
		return errors.New("config path is empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}
	Set(c)
	return nil
}

func DefaultPath() string {
	if path := os.Getenv("ONCALL_CONFIG"); path != "" {
		return path
	}
	candidates := []string{
		filepath.Join("manifest", "config", "config.yaml"),
		filepath.Join("..", "manifest", "config", "config.yaml"),
		filepath.Join("..", "..", "manifest", "config", "config.yaml"),
		filepath.Join("etc", "oncall-api.yaml"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return candidates[0]
}
