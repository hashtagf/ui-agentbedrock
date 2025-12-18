package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ui-agentbedrock/backend/internal/models"
	"github.com/ui-agentbedrock/backend/internal/repository"
	"github.com/ui-agentbedrock/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// MaxTokenEstimate is the estimated max tokens before auto-summarization
	// Claude 3 models (Sonnet, Haiku) have 200k context window
	// We use 80% threshold (160k) but account for agent call overhead (~40k)
	// So we trigger at ~120k tokens (80% of usable 150k context)
	MaxTokenEstimate = 120000
	// KeepRecentMessages is the number of recent messages to keep after summarization
	KeepRecentMessages = 4
)

type ChatHandler struct {
	agentService     *services.AgentService
	sessionService   *services.SessionService
	summarizeService *services.SummarizeService
	documentRepo     *repository.DocumentRepository
}

func NewChatHandler(agentService *services.AgentService, sessionService *services.SessionService, summarizeService *services.SummarizeService, documentRepo *repository.DocumentRepository) *ChatHandler {
	return &ChatHandler{
		agentService:     agentService,
		sessionService:   sessionService,
		summarizeService: summarizeService,
		documentRepo:     documentRepo,
	}
}

func (h *ChatHandler) StreamChat(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// Get existing messages to check token count
	_, messages, err := h.sessionService.GetSession(ctx, req.SessionID)
	if err != nil {
		log.Printf("Warning: Could not get session messages: %v", err)
	}

	// Get the AgentBedrock session ID (separate from MongoDB session ID)
	agentSessionID, summaryContext, err := h.sessionService.GetAgentSessionID(ctx, req.SessionID)
	if err != nil {
		log.Printf("Warning: Could not get agent session ID: %v", err)
		agentSessionID = req.SessionID // Fallback to MongoDB ID
	}

	// Check if we need to auto-summarize
	estimatedTokens := services.EstimateTokens(messages)
	summarized := false

	if estimatedTokens > MaxTokenEstimate && len(messages) > KeepRecentMessages {
		log.Printf("Auto-summarizing conversation (estimated %d tokens)", estimatedTokens)

		// Get messages to summarize (all except recent)
		messagesToSummarize := messages[:len(messages)-KeepRecentMessages]

		// Generate summary
		summary, err := h.summarizeService.SummarizeConversation(ctx, messagesToSummarize)
		if err != nil {
			log.Printf("Warning: Failed to summarize: %v", err)
		} else {
			// Save summary and clear old messages in MongoDB
			_, err = h.sessionService.SummarizeAndClearOld(ctx, req.SessionID, summary, int64(KeepRecentMessages))
			if err != nil {
				log.Printf("Warning: Failed to save summary: %v", err)
			} else {
				// CRITICAL: Rotate AgentBedrock session to reset its internal history
				newAgentSessionID, err := h.sessionService.RotateAgentSession(ctx, req.SessionID, summary)
				if err != nil {
					log.Printf("Warning: Failed to rotate agent session: %v", err)
				} else {
					agentSessionID = newAgentSessionID
					summaryContext = summary
					summarized = true
					log.Printf("Conversation summarized and agent session rotated: %s", agentSessionID)
				}
			}
		}
	}

	// Convert document IDs to ObjectIDs for saving
	docObjectIDs := make([]primitive.ObjectID, 0, len(req.DocumentIDs))
	for _, docIDStr := range req.DocumentIDs {
		docID, err := primitive.ObjectIDFromHex(docIDStr)
		if err == nil {
			docObjectIDs = append(docObjectIDs, docID)
		}
	}

	// Save user message with document IDs
	_, err = h.sessionService.SaveMessageWithDocuments(ctx, req.SessionID, "user", req.Message, docObjectIDs, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// Get document content if document IDs are provided
	documentContext := ""
	if len(req.DocumentIDs) > 0 {
		docIDs := make([]primitive.ObjectID, 0, len(req.DocumentIDs))
		for _, docIDStr := range req.DocumentIDs {
			docID, err := primitive.ObjectIDFromHex(docIDStr)
			if err != nil {
				log.Printf("Warning: Invalid document ID: %s", docIDStr)
				continue
			}
			docIDs = append(docIDs, docID)
		}

		if len(docIDs) > 0 {
			documents, err := h.documentRepo.GetDocumentsByIDs(ctx, docIDs)
			if err != nil {
				log.Printf("Warning: Failed to get documents: %v", err)
			} else {
				// Combine all document content
				var docContents []string
				for _, doc := range documents {
					if doc.Content != "" {
						docContents = append(docContents, fmt.Sprintf("[Document: %s]\n%s", doc.Filename, doc.Content))
					}
				}
				if len(docContents) > 0 {
					documentContext = strings.Join(docContents, "\n\n")
				}
			}
		}
	}

	// Prepare the message to send to AgentBedrock
	messageToSend := req.Message
	if documentContext != "" {
		messageToSend = fmt.Sprintf("[Document Context]\n%s\n\n[User Message]\n%s", documentContext, req.Message)
	}
	if summaryContext != "" {
		// Prepend summary context for the new AgentBedrock session
		// Keep sending summary context until we have enough new messages to replace it
		// Only clear summary context when we have accumulated enough new conversation (e.g., 10+ messages)
		messageCount, _ := h.sessionService.GetMessageCount(ctx, req.SessionID)
		messageToSend = fmt.Sprintf("[Previous Conversation Context]\n%s\n\n[Current Message]\n%s", summaryContext, req.Message)

		// Only clear summary context after accumulating enough new messages (10+ messages after summarization)
		// This ensures AI maintains context from the summary
		if messageCount >= 10 {
			h.sessionService.ClearSummaryContext(ctx, req.SessionID)
			log.Printf("Cleared summary context after accumulating %d messages", messageCount)
		} else {
			log.Printf("Applied summary context to agent session (message count: %d)", messageCount)
		}
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("X-Accel-Buffering", "no")

	// Create a channel for client disconnect detection
	clientGone := c.Request.Context().Done()
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// Stream callback
	callback := func(event models.SSEEvent) error {
		select {
		case <-clientGone:
			return fmt.Errorf("client disconnected")
		default:
			data, _ := json.Marshal(event.Data)
			_, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event.Event, string(data))
			if err != nil {
				return err
			}
			flusher.Flush()
			return nil
		}
	}

	// Notify client if summarization happened
	if summarized {
		summarizeData, _ := json.Marshal(map[string]interface{}{
			"message":        "Conversation history was automatically summarized to reduce context length",
			"newSessionId":   agentSessionID,
			"sessionRotated": true,
		})
		fmt.Fprintf(c.Writer, "event: summarized\ndata: %s\n\n", string(summarizeData))
		flusher.Flush()
	}

	// Invoke agent with streaming - use AgentBedrock session ID, not MongoDB ID
	trace, content, err := h.agentService.InvokeAgentStream(c.Request.Context(), agentSessionID, messageToSend, callback)

	// Save assistant message
	var assistantMessage *models.Message
	if content != "" {
		assistantMessage, _ = h.sessionService.SaveMessage(c.Request.Context(), req.SessionID, "assistant", content, trace)
	}

	// Send done event
	messageID := ""
	if assistantMessage != nil {
		messageID = assistantMessage.ID.Hex()
	}

	doneData, _ := json.Marshal(models.DoneEvent{MessageID: messageID})
	fmt.Fprintf(c.Writer, "event: done\ndata: %s\n\n", string(doneData))
	flusher.Flush()

	if err != nil && err != io.EOF {
		// Error already sent via callback
		return
	}
}
