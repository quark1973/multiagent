package retriever

import (
	"context"
	"github.com/quark1973/multiagent/internal/ai/embedder"
	"github.com/quark1973/multiagent/utility/client"
	"github.com/quark1973/multiagent/utility/common"

	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/components/retriever"
)

func NewMilvusRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		return nil, err
	}
	eb, err := embedder.DoubaoEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	r, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      cli,
		Collection:  common.MilvusCollectionName,
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:      1,
		Embedding: eb,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}
