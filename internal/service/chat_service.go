package service

import (
	"context"
	"errors"
	"github.com/quark1973/multiagent/internal/ai/agent/chat_pipeline"
	"github.com/quark1973/multiagent/utility/log_call_back"
	"github.com/quark1973/multiagent/utility/mem"
	"io"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type ChatService struct{}

func NewChatService() *ChatService {
	return &ChatService{}
}

func (s *ChatService) Chat(ctx context.Context, id string, question string) (string, error) {
	userMessage := &chat_pipeline.UserMessage{
		ID:      id,
		Query:   question,
		History: mem.GetSimpleMemory(id).GetMessages(),
	}

	runner, err := chat_pipeline.BuildChatAgent(ctx)
	if err != nil {
		return "", err
	}

	out, err := runner.Invoke(ctx, userMessage, compose.WithCallbacks(log_call_back.LogCallback(nil)))
	if err != nil {
		return "", err
	}
	mem.GetSimpleMemory(id).SetMessages(schema.UserMessage(question))
	mem.GetSimpleMemory(id).SetMessages(schema.SystemMessage(out.Content))
	return out.Content, nil
}

func (s *ChatService) Stream(ctx context.Context, id string, question string, onChunk func(string) error) error {
	userMessage := &chat_pipeline.UserMessage{
		ID:      id,
		Query:   question,
		History: mem.GetSimpleMemory(id).GetMessages(),
	}

	runner, err := chat_pipeline.BuildChatAgent(ctx)
	if err != nil {
		return err
	}
	stream, err := runner.Stream(ctx, userMessage, compose.WithCallbacks(log_call_back.LogCallback(nil)))
	if err != nil {
		return err
	}
	defer stream.Close()

	var fullResponse string
	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			if fullResponse != "" {
				mem.GetSimpleMemory(id).SetMessages(schema.UserMessage(question))
				mem.GetSimpleMemory(id).SetMessages(schema.SystemMessage(fullResponse))
			}
			return nil
		}
		if err != nil {
			return err
		}
		fullResponse += chunk.Content
		if onChunk != nil {
			if err := onChunk(chunk.Content); err != nil {
				return err
			}
		}
	}
}
