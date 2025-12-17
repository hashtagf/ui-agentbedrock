package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID primitive.ObjectID `bson:"session_id" json:"sessionId"`
	Role      string             `bson:"role" json:"role"` // "user" | "assistant"
	Content   string             `bson:"content" json:"content"`
	Trace     *Trace             `bson:"trace,omitempty" json:"trace,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

type ChatRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Message   string `json:"message" binding:"required"`
}
