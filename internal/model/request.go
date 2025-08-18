package model

// ChatMessage represents a chat message
type ChatMessage struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// ChatCompletionRequest represents the incoming request
type ChatCompletionRequest struct {
	Messages    []ChatMessage `json:"messages" binding:"required"`
	Model       string        `json:"model" binding:"required"`
	Temperature *float64      `json:"temperature,omitempty"`
	TopP        *float64      `json:"topP,omitempty"`
	Stream      *bool         `json:"stream,omitempty"`
}

// AkashChatRequest represents the request to Akash API
type AkashChatRequest struct {
	ID          string        `json:"id"`
	Messages    []ChatMessage `json:"messages"`
	Model       string        `json:"model"`
	System      string        `json:"system"`
	Temperature float64       `json:"temperature"`
	TopP        float64       `json:"topP"`
	Context     []interface{} `json:"context"`
}

// ImageStatusResponse represents the image generation status response
type ImageStatusResponse struct {
	JobID         string  `json:"job_id"`
	WorkerName    string  `json:"worker_name"`
	WorkerCity    string  `json:"worker_city"`
	WorkerCountry string  `json:"worker_country"`
	Status        string  `json:"status"`
	Result        string  `json:"result"`
	WorkerGPU     string  `json:"worker_gpu"`
	ElapsedTime   float64 `json:"elapsed_time"`
	QueuePosition int     `json:"queue_position"`
}