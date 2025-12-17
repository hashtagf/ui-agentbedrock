package services

import (
	"context"

	"github.com/ui-agentbedrock/backend/internal/models"
	"github.com/ui-agentbedrock/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionService struct {
	repo *repository.SessionRepository
}

func NewSessionService(repo *repository.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) CreateSession(ctx context.Context, title string) (*models.Session, error) {
	if title == "" {
		title = "New Chat"
	}

	session := &models.Session{
		Title: title,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) GetSessions(ctx context.Context) ([]models.Session, error) {
	return s.repo.GetSessions(ctx)
}

func (s *SessionService) GetSession(ctx context.Context, id string) (*models.Session, []models.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil, err
	}

	session, err := s.repo.GetSession(ctx, objectID)
	if err != nil {
		return nil, nil, err
	}

	messages, err := s.repo.GetMessages(ctx, objectID)
	if err != nil {
		return nil, nil, err
	}

	return session, messages, nil
}

func (s *SessionService) UpdateSession(ctx context.Context, id string, title string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return s.repo.UpdateSession(ctx, objectID, title)
}

func (s *SessionService) DeleteSession(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteSession(ctx, objectID)
}

func (s *SessionService) SaveMessage(ctx context.Context, sessionID string, role, content string, trace *models.Trace) (*models.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, err
	}

	message := &models.Message{
		SessionID: objectID,
		Role:      role,
		Content:   content,
		Trace:     trace,
	}

	if err := s.repo.SaveMessage(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *SessionService) UpdateMessageTrace(ctx context.Context, messageID string, trace *models.Trace) error {
	objectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return err
	}

	return s.repo.UpdateMessageTrace(ctx, objectID, trace)
}

// ClearMessages clears all messages from a session
func (s *SessionService) ClearMessages(ctx context.Context, sessionID string) error {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return err
	}

	return s.repo.ClearMessages(ctx, objectID)
}

// GetMessageCount returns the number of messages in a session
func (s *SessionService) GetMessageCount(ctx context.Context, sessionID string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return 0, err
	}

	return s.repo.GetMessageCount(ctx, objectID)
}

// GetRecentMessages gets the N most recent messages
func (s *SessionService) GetRecentMessages(ctx context.Context, sessionID string, limit int64) ([]models.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetRecentMessages(ctx, objectID, limit)
}

// SummarizeAndClearOld summarizes old messages and keeps only recent ones
func (s *SessionService) SummarizeAndClearOld(ctx context.Context, sessionID string, summary string, keepRecent int64) (*models.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, err
	}

	// Delete old messages
	if err := s.repo.DeleteOldMessages(ctx, objectID, keepRecent); err != nil {
		return nil, err
	}

	// Save summary as a system message
	summaryMessage := &models.Message{
		SessionID: objectID,
		Role:      "system",
		Content:   "[Conversation Summary]\n" + summary,
	}

	if err := s.repo.SaveMessage(ctx, summaryMessage); err != nil {
		return nil, err
	}

	return summaryMessage, nil
}
