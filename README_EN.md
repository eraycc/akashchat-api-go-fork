# akashchat-api-go

A Go-based REST API service that provides a proxy interface to Akash Chat API, supporting both text and image generation models.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-AGPL_v3-green.svg)](LICENSE)

[中文文档](README.md) | [English](README_EN.md)

## Features

- **Unified API Interface**: Compatible with OpenAI ChatGPT API format
- **Text Generation**: Support for various text generation models, including streaming support.
- **Image Generation**: Support for AkashGen image generation model
- **Session Management**: Automatic session token caching and refresh
- **Error Handling**: Comprehensive error handling and validation
- **Configurable**: Environment-based configuration
- **Docker Support**: Ready-to-use Docker configuration

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/006lp/akashchat-api-go.git
cd akashchat-api-go
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run cmd/server/main.go
```

The server will start on `localhost:16571` by default.

### Using Docker

1. Build the Docker image:
```bash
docker build -t akashchat-api-go .
```

2. Run the container:
```bash
docker run -p 16571:16571 akashchat-api-go
```

## API Usage

### Text Generation

Send a POST request to `/v1/chat/completions`:

```bash
curl -X POST http://localhost:16571/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "Hello, how are you?"
      }
    ],
    "model": "Meta-Llama-3-3-70B-Instruct",
    "temperature": 0.85,
    "topP": 1.0
  }'
```

**Response:**
```json
{
  "choices": [
    {
      "finish_reason": "stop",
      "index": 0,
      "message": {
        "content": "Hello! I'm doing well, thank you for asking...",
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

### Image Generation

Use the `AkashGen` model for image generation:

```bash
curl -X POST http://localhost:16571/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "a cute anime girl with blue eyes"
      }
    ],
    "model": "AkashGen",
    "temperature": 0.85
  }'
```

**Response:**
```json
{
  "code": 200,
  "data": {
    "model": "AkashGen",
    "jobId": "727ef62f-76c9-45b8-9637-dc461590fe49",
    "prompt": "a cute anime girl with pastel-colored hair and sparkling blue eyes...",
    "pic": "https://chat.akash.network/api/image/job_727ef62f_00001_.webp"
  }
}
```

### Get Model List

Get a list of all available models:

```bash
curl http://localhost:16571/v1/models
```

**Response:**
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

### Health Check

Check if the service is running:

```bash
curl http://localhost:16571/health
```

## Configuration

The application can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_ADDRESS` | `localhost:16571` | Server address and port |
| `AKASH_BASE_URL` | `https://chat.akash.network` | Akash Chat API base URL |

Example:
```bash
export SERVER_ADDRESS="0.0.0.0:8080"
export AKASH_BASE_URL="https://chat.akash.network"
go run cmd/server/main.go
```

## Project Structure

```
akashchat-api-go/
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── handler/         # HTTP request handlers
│   ├── model/          # Data models
│   ├── service/        # Business logic
│   └── utils/          # Utility functions
├── pkg/                # Public packages
│   └── client/         # HTTP client wrapper
├── Dockerfile          # Docker configuration
└── README.md           # This file
```

## API Parameters

### Request Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `messages` | Array | Yes | - | Array of message objects |
| `model` | String | Yes | - | Model name (e.g., "Meta-Llama-3-3-70B-Instruct", "AkashGen") |
| `temperature` | Float | No | 0.85 | Sampling temperature (0.0-2.0) |
| `topP` | Float | No | 1.0 | Top-p sampling parameter |
| `stream` | Boolean | No | false | Enable streaming response |

### Message Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `role` | String | Yes | Message role ("user", "assistant", "system") |
| `content` | String | Yes | Message content |

## Error Handling

The API returns standardized error responses:

```json
{
  "code": 500,
  "data": {
    "msg": "Error Model."
  }
}
```

Common error codes:
- `400`: Bad Request (invalid JSON or missing required fields)
- `500`: Internal Server Error (invalid model, API errors)

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o bin/akashchat-api-go cmd/server/main.go
```

### Code Structure

The project follows standard Go project layout:

- **cmd/**: Main applications for this project
- **internal/**: Private application and library code
- **pkg/**: Library code that's ok to use by external applications

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the AGPL v3 License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin) - HTTP web framework
- [Akash Network](https://akash.network/) - Decentralized cloud compute platform

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/006lp/akashchat-api-go/issues) on GitHub.