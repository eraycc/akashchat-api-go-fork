# akashchat-api-go

基于 Go 语言开发的 REST API 服务，为 Akash Chat API 提供代理接口，支持文本和图像生成模型。

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-AGPL_v3-green.svg)](LICENSE)

[中文文档](README.md) | [English](README_EN.md)

## 功能特性

- **统一 API 接口**: 兼容 OpenAI ChatGPT API 格式
- **文本生成**: 支持多种文本生成模型，并支持流式响应
- **图像生成**: 支持 AkashGen 图像生成模型
- **会话管理**: 自动会话令牌缓存和刷新
- **错误处理**: 全面的错误处理和验证
- **可配置**: 基于环境变量的配置
- **Docker 支持**: 提供即用的 Docker 配置

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- Git

### 安装

1. 克隆仓库:
```bash
git clone https://github.com/006lp/akashchat-api-go.git
cd akashchat-api-go
```

2. 安装依赖:
```bash
go mod tidy
```

3. 运行应用:
```bash
go run cmd/server/main.go
```

服务器默认在 `localhost:16571` 启动。

### 使用 Docker

1. 构建 Docker 镜像:
```bash
docker build -t akashchat-api-go .
```

2. 运行容器:
```bash
docker run -p 16571:16571 akashchat-api-go
```

## API 使用

### 文本生成

向 `/v1/chat/completions` 发送 POST 请求:

```bash
curl -X POST http://localhost:16571/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "你好，你怎么样？"
      }
    ],
    "model": "Meta-Llama-3-3-70B-Instruct",
    "temperature": 0.85,
    "topP": 1.0
  }'
```

**响应:**
```json
{
  "choices": [
    {
      "finish_reason": "stop",
      "index": 0,
      "message": {
        "content": "你好！我很好，谢谢你的询问...",
        "role": "assistant"
      }
    }
  ],
  "created": 1755506652,
  "id": "chatcmpl-1755506652.79333",
  "model": "Meta-Llama-3-3-70B-Instruct",
  "object": "chat.completion",
  "usage": {
    "completion_tokens": 0,
    "prompt_tokens": 0,
    "total_tokens": 0
  }
}
```

### 图像生成

使用 `AkashGen` 模型进行图像生成:

```bash
curl -X POST http://localhost:16571/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "一个蓝眼睛的可爱动漫女孩"
      }
    ],
    "model": "AkashGen",
    "temperature": 0.85
  }'
```

**响应:**
```json
{
  "code": 200,
  "data": {
    "model": "AkashGen",
    "jobId": "727ef62f-76c9-45b8-9637-dc461590fe49",
    "prompt": "一个蓝眼睛的可爱动漫女孩，有着彩色头发和闪亮的蓝眼睛...",
    "pic": "https://chat.akash.network/api/image/job_727ef62f_00001_.webp"
  }
}
```

### 获取模型列表

获取所有可用模型的列表：

```bash
curl http://localhost:16571/v1/models
```

**响应:**
```json
[
  {
    "id": "openai-gpt-oss-120b",
    "name": "GPT-OSS-120B",
    "description": "Efficient reasoning model with 117B parameters (5.1B active)",
    "temperature": 0.6,
    "top_p": 0.95,
    "tokenLimit": 128000,
    "parameters": "117B (5.1B active)",
    "architecture": "Transformer with native MXFP4 quantization",
    "hf_repo": "openai/gpt-oss-120b",
    "aboutContent": "Experience **GPT-OSS-120B**, OpenAI's open-source reasoning model with 117B total parameters and 5.1B active parameters. Built with native MXFP4 quantization, this model is designed for powerful reasoning, agentic tasks, and versatile developer use cases.\n\nGPT-OSS-120B features configurable reasoning levels (Low, Medium, High) and supports advanced capabilities like tool use, web browsing, and function calling. Optimized to run efficiently on a single H100 GPU while delivering high-quality reasoning performance.",
    "infoContent": "\n* ⚡ Open-source reasoning model with configurable reasoning levels\n* 🧠 117B parameters (5.1B active) with native MXFP4 quantization\n* 🌐 Decentralized hosting for lower costs & full control\n* 🔍 Optimized for reasoning, agentic tasks, and tool use",
    "thumbnailId": "llama-3",
    "deployUrl": "https://console.akash.network/templates/akash-network-awesome-akash-openai-gpt-oss-120b",
    "available": true
  }
]
```

### 健康检查

检查服务是否正常运行:

```bash
curl http://localhost:16571/health
```

## 配置

应用程序可以通过环境变量进行配置:

| 变量 | 默认值 | 描述 |
|------|--------|------|
| `SERVER_ADDRESS` | `localhost:16571` | 服务器地址和端口 |
| `AKASH_BASE_URL` | `https://chat.akash.network` | Akash Chat API 基础 URL |

示例:
```bash
export SERVER_ADDRESS="0.0.0.0:8080"
export AKASH_BASE_URL="https://chat.akash.network"
go run cmd/server/main.go
```

## 项目结构

```
akashchat-api-go/
├── cmd/server/          # 应用程序入口
├── internal/            # 私有应用程序代码
│   ├── config/          # 配置管理
│   ├── handler/         # HTTP 请求处理器
│   ├── model/          # 数据模型
│   ├── service/        # 业务逻辑
│   └── utils/          # 工具函数
├── pkg/                # 公共包
│   └── client/         # HTTP 客户端封装
├── Dockerfile          # Docker 配置
└── README.md           # 说明文档
```

## API 参数

### 请求参数

| 参数 | 类型 | 必需 | 默认值 | 描述 |
|------|------|------|--------|------|
| `messages` | 数组 | 是 | - | 消息对象数组 |
| `model` | 字符串 | 是 | - | 模型名称（例如："Meta-Llama-3-3-70B-Instruct"、"AkashGen"） |
| `temperature` | 浮点数 | 否 | 0.85 | 采样温度（0.0-2.0） |
| `topP` | 浮点数 | 否 | 1.0 | Top-p 采样参数 |
| `stream` | 布尔值 | 否 | false | 是否启用流式响应 |

### 消息对象

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `role` | 字符串 | 是 | 消息角色（"user"、"assistant"、"system"） |
| `content` | 字符串 | 是 | 消息内容 |

## 错误处理

API 返回标准化的错误响应:

```json
{
  "code": 500,
  "data": {
    "msg": "Error Model."
  }
}
```

常见错误代码:
- `400`: 请求错误（无效 JSON 或缺少必需字段）
- `500`: 内部服务器错误（无效模型、API 错误）

## 开发

### 运行测试

```bash
go test ./...
```

### 构建

```bash
go build -o bin/akashchat-api-go cmd/server/main.go
```

### 代码结构

项目遵循标准的 Go 项目布局:

- **cmd/**: 项目的主要应用程序
- **internal/**: 私有应用程序和库代码
- **pkg/**: 可以被外部应用程序使用的库代码

## 贡献

1. Fork 仓库
2. 创建你的功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的更改 (`git commit -m '添加一些很棒的功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 AGPL v3 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 致谢

- [Gin Web Framework](https://github.com/gin-gonic/gin) - HTTP Web 框架
- [Akash Network](https://akash.network/) - 去中心化云计算平台

## 支持

如果你遇到任何问题或有疑问，请在 GitHub 上[创建 issue](https://github.com/006lp/akashchat-api-go/issues)。

## 部署说明

### 本地开发

1. 确保已安装 Go 1.21+
2. 克隆项目并进入目录
3. 运行 `go mod tidy` 安装依赖
4. 运行 `go run cmd/server/main.go` 启动服务

### 生产部署

推荐使用 Docker 进行部署：

```bash
# 构建镜像
docker build -t akashchat-api-go:latest .

# 运行容器
docker run -d \
  --name akashchat-api \
  -p 16571:16571 \
  -e SERVER_ADDRESS="0.0.0.0:16571" \
  akashchat-api-go:latest
```

### 注意事项

- 会话令牌有效期为 1 小时，系统会自动处理刷新
- 图像生成请求可能需要较长时间，系统会自动轮询直到完成
- 建议在生产环境中设置适当的请求超时和限流策略