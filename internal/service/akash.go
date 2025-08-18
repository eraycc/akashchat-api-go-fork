package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/006lp/akashchat-api-go/internal/model"
	"github.com/006lp/akashchat-api-go/internal/utils"
	"github.com/006lp/akashchat-api-go/pkg/client"
)

// AkashService handles communication with Akash API
type AkashService struct {
	httpClient *client.HTTPClient
}

// NewAkashService creates a new AkashService instance
func NewAkashService() *AkashService {
	return &AkashService{
		httpClient: client.NewHTTPClient(),
	}
}

// ProcessImageGeneration handles image generation requests
func (a *AkashService) ProcessImageGeneration(req model.ChatCompletionRequest, sessionToken string, temperature, topP float64) (*model.ImageGenerationData, error) {
	// Create Akash chat request
	akashReq := model.AkashChatRequest{
		ID:          utils.GenerateRandomID(16),
		Messages:    req.Messages,
		Model:       req.Model,
		System:      getSystemPrompt(),
		Temperature: temperature,
		TopP:        topP,
		Context:     []interface{}{},
	}

	// Send chat request
	respText, err := a.sendChatRequest(akashReq, sessionToken)
	if err != nil {
		return nil, err
	}

	// Check for error response
	if strings.Contains(respText, "error") && strings.Contains(respText, "Invalid model name") {
		return nil, fmt.Errorf("invalid model")
	}

	// Extract jobId and prompt from response
	jobID, prompt, err := a.extractImageGenerationInfo(respText)
	if err != nil {
		return nil, fmt.Errorf("failed to extract image generation info: %w", err)
	}

	// Poll for image completion
	imageURL, err := a.pollImageStatus(jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get image result: %w", err)
	}

	return &model.ImageGenerationData{
		Model:  req.Model,
		JobID:  jobID,
		Prompt: prompt,
		Pic:    imageURL,
	}, nil
}

// ProcessTextGeneration handles text generation requests
func (a *AkashService) ProcessTextGeneration(req model.ChatCompletionRequest, sessionToken string, temperature, topP float64) (*model.OpenAIChatCompletion, error) {
	// Create Akash chat request
	akashReq := model.AkashChatRequest{
		ID:          utils.GenerateRandomID(16),
		Messages:    req.Messages,
		Model:       req.Model,
		System:      getSystemPrompt(),
		Temperature: temperature,
		TopP:        topP,
		Context:     []interface{}{},
	}

	// Send chat request
	respText, err := a.sendChatRequest(akashReq, sessionToken)
	if err != nil {
		return nil, err
	}

	// Check for error response
	if strings.Contains(respText, "error") && strings.Contains(respText, "Invalid model name") {
		return nil, fmt.Errorf("invalid model")
	}

	// Extract and format the response
	openAIResp := a.extractTextGenerationInfo(respText, req.Model)

	return openAIResp, nil
}

// sendChatRequest sends a request to Akash chat API
func (a *AkashService) sendChatRequest(req model.AkashChatRequest, sessionToken string) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	headers := map[string]string{
		"Referer":      "https://chat.akash.network/",
		"Cookie":       sessionToken,
		"Accept":       "*/*",
		"Content-Type": "application/json",
	}

	resp, err := a.httpClient.Post("https://chat.akash.network/api/chat/", bytes.NewBuffer(jsonData), headers)
	if err != nil {
		return "", fmt.Errorf("failed to send chat request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(respBody), nil
}

// extractImageGenerationInfo extracts jobId and prompt from image generation response
func (a *AkashService) extractImageGenerationInfo(respText string) (string, string, error) {
	// Extract jobId using regex
	jobIDRegex := regexp.MustCompile(`jobId='([^']+)'`)
	jobIDMatch := jobIDRegex.FindStringSubmatch(respText)
	if len(jobIDMatch) < 2 {
		return "", "", fmt.Errorf("jobId not found in response")
	}
	jobID := jobIDMatch[1]

	// Extract prompt using regex
	promptRegex := regexp.MustCompile(`prompt='([^']+)'`)
	promptMatch := promptRegex.FindStringSubmatch(respText)
	if len(promptMatch) < 2 {
		return "", "", fmt.Errorf("prompt not found in response")
	}
	prompt := promptMatch[1]

	return jobID, prompt, nil
}

// pollImageStatus polls the image status until completion
func (a *AkashService) pollImageStatus(jobID string) (string, error) {
	maxAttempts := 60 // Maximum 1 minute polling
	
	for i := 0; i < maxAttempts; i++ {
		url := fmt.Sprintf("https://chat.akash.network/api/image-status?ids=%s", jobID)
		
		resp, err := a.httpClient.Get(url, nil)
		if err != nil {
			return "", fmt.Errorf("failed to check image status: %w", err)
		}

		var statusResp []model.ImageStatusResponse
		if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
			resp.Body.Close()
			return "", fmt.Errorf("failed to decode image status response: %w", err)
		}
		resp.Body.Close()

		if len(statusResp) == 0 {
			return "", fmt.Errorf("empty status response")
		}

		status := statusResp[0]
		if status.Status == "succeeded" {
			return "https://chat.akash.network" + status.Result, nil
		}

		if status.Status == "failed" {
			return "", fmt.Errorf("image generation failed")
		}

		// Wait 1 second before next poll
		time.Sleep(1 * time.Second)
	}

	return "", fmt.Errorf("image generation timed out")
}

// extractTextGenerationInfo extracts text generation information from response and formats it as OpenAI's chat completion.
func (a *AkashService) extractTextGenerationInfo(respText string, modelName string) *model.OpenAIChatCompletion {
	var messageID string
	var allContent strings.Builder
	var finishReason string

	lines := strings.Split(respText, "\n")
	var contentStarted bool

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Extract messageId
		if strings.HasPrefix(line, "f:{\"messageId\":") {
			msgIDRegex := regexp.MustCompile(`"messageId":"([^"]+)"`)
			match := msgIDRegex.FindStringSubmatch(line)
			if len(match) > 1 {
				messageID = match[1]
			}
			contentStarted = true
			continue
		}

		// Extract finishReason
		if strings.HasPrefix(line, "e:{\"finishReason\":") {
			reasonRegex := regexp.MustCompile(`"finishReason":"([^"]+)"`)
			match := reasonRegex.FindStringSubmatch(line)
			if len(match) > 1 {
				finishReason = match[1]
			}
			break
		}

		// Collect content lines
		if contentStarted && strings.HasPrefix(line, "0:\"") {
			content := line[3:]
			content = strings.TrimSuffix(content, "\"")
			content = strings.ReplaceAll(content, "\\n", "\n")
			content = strings.ReplaceAll(content, "\\\"", "\"")
			allContent.WriteString(content)
		}
	}

	fullContent := allContent.String()

	// Create OpenAI format response
	return &model.OpenAIChatCompletion{
		ID:      "chatcmpl-" + messageID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   modelName,
		Choices: []model.Choice{
			{
				Index: 0,
				Message: model.Message{
					Role:    "assistant",
					Content: fullContent,
				},
				FinishReason: finishReason,
			},
		},
		Usage: model.Usage{
			PromptTokens:     0, // Placeholder, as Akash does not provide this
			CompletionTokens: 0, // Placeholder
			TotalTokens:      0, // Placeholder
		},
	}
}

// getSystemPrompt returns the system prompt for Akash
func getSystemPrompt() string {
	return "You are a skilled conversationalist who adapts naturally to what users need. Your responses match the situation—whether someone wants deep analysis, casual chat, emotional support, creative collaboration, or just needs to vent.\nCore Approach\n\nRead between the lines to understand what people actually want\nMatch their energy and conversational style\nShift seamlessly between modes: analytical, empathetic, humorous, creative, or practical\nWhen people need to be heard, focus on listening rather than fixing\nFor substantive topics, provide thorough, well-organized insights that aid decision-making\n\nCommunication Style\n\nSound natural and authentic, never templated or robotic\nAvoid unnecessary politeness policing or inclusion reminders\nWrite in requested voices, styles, or perspectives when asked\nAdapt tone appropriately—you can be direct, irreverent, or even rude when specifically prompted to do so\n\nInteraction Philosophy\n\nSometimes the best help is simply being present and understanding\nDon't over-optimize for helpfulness when someone just wants connection\nTrust that users know what they're looking for and deliver accordingly\nProvide depth and insight for complex topics while keeping casual conversations light"
}