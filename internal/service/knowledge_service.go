package service

import (
	"context"
	"fmt"
	"github.com/quark1973/multiagent/internal/ai/agent/knowledge_index_pipeline"
	loader2 "github.com/quark1973/multiagent/internal/ai/loader"
	"github.com/quark1973/multiagent/utility/client"
	"github.com/quark1973/multiagent/utility/common"
	"github.com/quark1973/multiagent/utility/log_call_back"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
)

type UploadedFile struct {
	FileName string
	FilePath string
	FileSize int64
}

type KnowledgeService struct{}

func NewKnowledgeService() *KnowledgeService {
	return &KnowledgeService{}
}

func (s *KnowledgeService) UploadAndIndex(ctx context.Context, r *http.Request) (*UploadedFile, error) {
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		return nil, err
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("missing upload file: %w", err)
	}
	defer file.Close()

	if err := os.MkdirAll(common.FileDir, 0755); err != nil {
		return nil, fmt.Errorf("create file dir failed: %w", err)
	}

	fileName := filepath.Base(header.Filename)
	savePath := filepath.Join(common.FileDir, fileName)
	dst, err := os.Create(savePath)
	if err != nil {
		return nil, fmt.Errorf("create upload file failed: %w", err)
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("save upload file failed: %w", err)
	}

	if err := s.buildIntoIndex(ctx, savePath); err != nil {
		return nil, err
	}

	return &UploadedFile{
		FileName: fileName,
		FilePath: savePath,
		FileSize: size,
	}, nil
}

func (s *KnowledgeService) buildIntoIndex(ctx context.Context, path string) error {
	runner, err := knowledge_index_pipeline.BuildKnowledgeIndexing(ctx)
	if err != nil {
		return err
	}
	loader, err := loader2.NewFileLoader(ctx)
	if err != nil {
		return err
	}
	docs, err := loader.Load(ctx, document.Source{URI: path})
	if err != nil {
		return err
	}
	if len(docs) == 0 {
		return fmt.Errorf("no document loaded from %s", path)
	}

	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		return err
	}
	source := fmt.Sprint(docs[0].MetaData["_source"])
	expr := fmt.Sprintf(`metadata["_source"] == "%s"`, source)
	queryResult, err := cli.Query(ctx, common.MilvusCollectionName, []string{}, expr, []string{"id"})
	if err != nil {
		return err
	}

	var idsToDelete []string
	for _, column := range queryResult {
		if column.Name() != "id" {
			continue
		}
		for i := 0; i < column.Len(); i++ {
			id, err := column.GetAsString(i)
			if err == nil {
				idsToDelete = append(idsToDelete, id)
			}
		}
	}
	if len(idsToDelete) > 0 {
		deleteExpr := fmt.Sprintf(`id in ["%s"]`, strings.Join(idsToDelete, `","`))
		if err := cli.Delete(ctx, common.MilvusCollectionName, "", deleteExpr); err != nil {
			return err
		}
	}

	_, err = runner.Invoke(ctx, document.Source{URI: path}, compose.WithCallbacks(log_call_back.LogCallback(nil)))
	return err
}
