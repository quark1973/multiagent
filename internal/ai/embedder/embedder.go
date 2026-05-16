package embedder

import (
	"context"
	"github.com/quark1973/multiagent/utility/appconfig"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/cloudwego/eino/components/embedding"
)

func DoubaoEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	cfg, err := appconfig.Get()
	if err != nil {
		return nil, err
	}
	modelCfg := cfg.DoubaoEmbeddingModel
	dim := 2048
	embedder, err := dashscope.NewEmbedder(ctx, &dashscope.EmbeddingConfig{
		Model:      modelCfg.Model,
		APIKey:     modelCfg.APIKey,
		Dimensions: &dim,
	})
	if err != nil {
		log.Printf("new embedder error: %v\n", err)
		return nil, err
	}
	return embedder, nil
}
