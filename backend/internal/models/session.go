package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title          string             `bson:"title" json:"title"`
	AgentSessionID string             `bson:"agent_session_id" json:"agentSessionId"`                    // Separate ID for AgentBedrock API
	SummaryContext string             `bson:"summary_context,omitempty" json:"summaryContext,omitempty"` // Context to pass on session rotation
	CreatedAt      time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updatedAt"`
}

type CreateSessionRequest struct {
	Title string `json:"title"`
}

type UpdateSessionRequest struct {
	Title string `json:"title"`
}
