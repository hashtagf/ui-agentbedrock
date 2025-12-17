package repository

import (
	"context"
	"time"

	"github.com/ui-agentbedrock/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SessionRepository struct {
	sessions *mongo.Collection
	messages *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) *SessionRepository {
	return &SessionRepository{
		sessions: db.Collection("sessions"),
		messages: db.Collection("messages"),
	}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *models.Session) error {
	session.ID = primitive.NewObjectID()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	_, err := r.sessions.InsertOne(ctx, session)
	return err
}

func (r *SessionRepository) GetSessions(ctx context.Context) ([]models.Session, error) {
	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})
	cursor, err := r.sessions.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.Session
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	if sessions == nil {
		sessions = []models.Session{}
	}
	return sessions, nil
}

func (r *SessionRepository) GetSession(ctx context.Context, id primitive.ObjectID) (*models.Session, error) {
	var session models.Session
	err := r.sessions.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) UpdateSession(ctx context.Context, id primitive.ObjectID, title string) error {
	_, err := r.sessions.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"title":      title,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

func (r *SessionRepository) DeleteSession(ctx context.Context, id primitive.ObjectID) error {
	// Delete all messages in the session
	_, err := r.messages.DeleteMany(ctx, bson.M{"session_id": id})
	if err != nil {
		return err
	}

	// Delete the session
	_, err = r.sessions.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *SessionRepository) GetMessages(ctx context.Context, sessionID primitive.ObjectID) ([]models.Message, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	cursor, err := r.messages.Find(ctx, bson.M{"session_id": sessionID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	if messages == nil {
		messages = []models.Message{}
	}
	return messages, nil
}

func (r *SessionRepository) SaveMessage(ctx context.Context, message *models.Message) error {
	message.ID = primitive.NewObjectID()
	message.CreatedAt = time.Now()

	_, err := r.messages.InsertOne(ctx, message)
	if err != nil {
		return err
	}

	// Update session's updated_at
	_, err = r.sessions.UpdateOne(
		ctx,
		bson.M{"_id": message.SessionID},
		bson.M{"$set": bson.M{"updated_at": time.Now()}},
	)
	return err
}

func (r *SessionRepository) UpdateMessageTrace(ctx context.Context, messageID primitive.ObjectID, trace *models.Trace) error {
	_, err := r.messages.UpdateOne(
		ctx,
		bson.M{"_id": messageID},
		bson.M{"$set": bson.M{"trace": trace}},
	)
	return err
}

// ClearMessages deletes all messages for a session
func (r *SessionRepository) ClearMessages(ctx context.Context, sessionID primitive.ObjectID) error {
	_, err := r.messages.DeleteMany(ctx, bson.M{"session_id": sessionID})
	return err
}

// GetMessageCount returns the number of messages in a session
func (r *SessionRepository) GetMessageCount(ctx context.Context, sessionID primitive.ObjectID) (int64, error) {
	return r.messages.CountDocuments(ctx, bson.M{"session_id": sessionID})
}

// GetRecentMessages gets the N most recent messages
func (r *SessionRepository) GetRecentMessages(ctx context.Context, sessionID primitive.ObjectID, limit int64) ([]models.Message, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit)
	cursor, err := r.messages.Find(ctx, bson.M{"session_id": sessionID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// DeleteOldMessages deletes messages except the most recent N
func (r *SessionRepository) DeleteOldMessages(ctx context.Context, sessionID primitive.ObjectID, keepRecent int64) error {
	// Get IDs of recent messages to keep
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(keepRecent)
	cursor, err := r.messages.Find(ctx, bson.M{"session_id": sessionID}, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var recentMessages []models.Message
	if err := cursor.All(ctx, &recentMessages); err != nil {
		return err
	}

	// Collect IDs to keep
	keepIDs := make([]primitive.ObjectID, len(recentMessages))
	for i, msg := range recentMessages {
		keepIDs[i] = msg.ID
	}

	// Delete all except recent
	_, err = r.messages.DeleteMany(ctx, bson.M{
		"session_id": sessionID,
		"_id":        bson.M{"$nin": keepIDs},
	})
	return err
}
