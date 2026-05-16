# MultiAgent OnCall Assistant

基于 go-zero + CloudWeGo Eino 的智能 OnCall Agent 服务，支持企业知识库问答、文档向量化索引、AIOps 告警分析，以及飞书机器人接入。

## 功能特性

- go-zero 分层服务结构：`handler / logic / svc / service / config`
- RAG 知识库：支持文档上传、切分、Embedding 和 Milvus 向量索引
- Agent 编排：普通问答使用 ReAct Agent，复杂告警分析使用 Plan-Execute-Replan
- 工具调用：接入日志 MCP、Prometheus 告警查询、MySQL 查询、内部文档检索等工具
- 飞书机器人：支持飞书群聊消息事件回调，并将 Agent 结果回传到原群聊
- SSE 流式问答：支持 `/api/chat_stream` 实时返回模型输出

## 项目结构

```text
.
├── oncall.go                 # go-zero 服务入口
├── oncall.api                # goctl API 定义
├── etc/oncall-api.yaml       # 服务配置
├── internal/
│   ├── handler               # HTTP handler
│   ├── logic                 # go-zero logic 层
│   ├── service               # Chat / Knowledge / AIOps / Feishu 业务服务
│   ├── svc                   # ServiceContext
│   ├── types                 # API DTO
│   └── ai                    # Eino Agent、RAG、模型、工具调用
├── utility/                  # 配置、Milvus、日志回调、会话记忆
└── docs/                     # 示例知识库文档
```

## 接口列表

| Method | Path | Description |
| --- | --- | --- |
| POST | `/api/chat` | 普通 Agent 问答 |
| POST | `/api/chat_stream` | SSE 流式问答 |
| POST | `/api/upload` | 上传文档并构建知识库索引 |
| POST | `/api/ai_ops` | 智能运维告警分析 |
| POST | `/api/feishu/event` | 飞书机器人事件回调 |

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 修改配置

编辑 `etc/oncall-api.yaml`，填入模型、Embedding、MCP、飞书应用等配置：

```yaml
ds_quick_chat_model:
  api_key: "your-model-api-key"
  base_url: "https://ark.cn-beijing.volces.com/api/v3"
  model: "deepseek-v3-1-terminus"

doubao_embedding_model:
  api_key: "your-dashscope-api-key"
  base_url: "https://dashscope.aliyuncs.com/compatible-mode/v1"
  model: "text-embedding-v4"

mcp_url: "https://mcp-api.tencent-cloud.com/sse/XXXX"

feishu:
  app_id: "your-feishu-app-id"
  app_secret: "your-feishu-app-secret"
  verify_token: "your-feishu-verify-token"
  encrypt_key: "your-feishu-encrypt-key"
```

### 3. 启动服务

```bash
go run . -f etc/oncall-api.yaml
```

默认端口为 `6872`。

## 飞书接入说明

在飞书开放平台创建机器人应用后，将事件订阅请求地址配置为：

```text
https://your-domain.com/api/feishu/event
```

消息路由规则：

- 普通文本消息：进入 ChatAgent 问答链路
- `/aiops` 或包含“告警分析”的消息：进入 AIOps 分析链路

机器人会先回传 `received, analyzing...`，任务完成后再把最终结果发回原群聊。

## 运行检查

```bash
go test ./...
```

## 注意事项

- 本项目需要可用的模型 API Key、Embedding API Key、Milvus、MCP 地址和飞书应用配置。
- `queryPrometheusAlerts` 当前保留了演示逻辑，如需真实查询 Prometheus，需要补充实际 Prometheus 地址与查询实现。
- 不要将真实 API Key、App Secret 或私有配置提交到公开仓库。
