package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/ui-agentbedrock/backend/internal/models"
)

type SummarizeService struct {
	bedrockClient *bedrockruntime.Client
	modelID       string
}

func NewSummarizeService(cfg aws.Config) *SummarizeService {
	return &SummarizeService{
		bedrockClient: bedrockruntime.NewFromConfig(cfg),
		modelID:       "anthropic.claude-3-haiku-20240307-v1:0", // Fast and cheap for summarization
	}
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version"`
	MaxTokens        int             `json:"max_tokens"`
	System           string          `json:"system,omitempty"`
	Messages         []claudeMessage `json:"messages"`
}

type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

// SummarizeConversation summarizes a list of messages into a concise summary
func (s *SummarizeService) SummarizeConversation(ctx context.Context, messages []models.Message) (string, error) {
	if len(messages) == 0 {
		return "", nil
	}

	// Build conversation text
	var conversationParts []string
	for _, msg := range messages {
		var role string
		switch msg.Role {
		case "user":
			role = "User"
		case "assistant":
			role = "AI"
		case "system":
			role = "System"
		default:
			role = msg.Role
		}
		conversationParts = append(conversationParts, fmt.Sprintf("%s: %s", role, msg.Content))
	}
	conversationText := strings.Join(conversationParts, "\n\n")

	// Create summarization prompt
	systemPrompt := `You are a conversation summarizer. Your task is to create a concise but comprehensive summary of the conversation history. 
Focus on:
- Key topics discussed
- Important decisions or conclusions reached
- Any pending questions or tasks
- Context that would be needed to continue the conversation

Keep the summary under 500 words. Be factual and objective.`

	userPrompt := fmt.Sprintf(`Please summarize the following conversation:

%s

Provide a concise summary that captures the essential context needed to continue this conversation.`, conversationText)

	// Build request
	requestBody := claudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        1024,
		System:           systemPrompt,
		Messages: []claudeMessage{
			{Role: "user", Content: userPrompt},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Invoke model
	output, err := s.bedrockClient.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(s.modelID),
		ContentType: aws.String("application/json"),
		Body:        bodyBytes,
	})
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	// Parse response
	var response claudeResponse
	if err := json.Unmarshal(output.Body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return response.Content[0].Text, nil
}

// EstimateTokens provides a rough estimate of token count
// Claude uses roughly 4 characters per token on average
func EstimateTokens(messages []models.Message) int {
	totalChars := 0
	for _, msg := range messages {
		totalChars += len(msg.Content)
	}
	return totalChars / 4
}
