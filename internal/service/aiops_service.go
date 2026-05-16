package service

import (
	"context"
	"errors"
	"github.com/quark1973/multiagent/internal/ai/agent/plan_execute_replan"
)

const defaultAIOpsPrompt = `
1. You are an intelligent service alert analysis assistant. First call query_prometheus_alerts to get all active alerts.
2. For each alert name, call query_internal_docs to get the matching handling SOP.
3. Analyze strictly based on internal documents. Do not invent information outside the retrieved docs.
4. For any time parameter, call get_current_time first and then pass the required timestamp.
5. For log investigation, call the log tool with the configured region and log topic.
6. Summarize the information into an alert operation analysis report with active alerts, root-cause analysis, handling steps, and conclusion.
`

type AIOpsService struct{}

func NewAIOpsService() *AIOpsService {
	return &AIOpsService{}
}

func (s *AIOpsService) Analyze(ctx context.Context) (string, []string, error) {
	resp, detail, err := plan_execute_replan.BuildPlanAgent(ctx, defaultAIOpsPrompt)
	if err != nil {
		return "", nil, err
	}
	if resp == "" {
		return "", nil, errors.New("empty aiops response")
	}
	return resp, detail, nil
}
