package models

import (
	"time"
)

type Trace struct {
	TraceID    string      `bson:"trace_id" json:"traceId"`
	AgentSteps []AgentStep `bson:"agent_steps" json:"agentSteps"`
	Error      *ErrorInfo  `bson:"error,omitempty" json:"error,omitempty"`
}

type AgentStep struct {
	StepIndex   int       `bson:"step_index" json:"stepIndex"`
	AgentName   string    `bson:"agent_name" json:"agentName"`
	AgentID     string    `bson:"agent_id,omitempty" json:"agentId,omitempty"`
	Type        string    `bson:"type" json:"type"` // "orchestration" | "pre_processing" | "post_processing" | "action" | "knowledge_base" | "collaborator"
	Action      string    `bson:"action" json:"action"`
	Status      string    `bson:"status" json:"status"` // "running" | "success" | "error"
	Rationale   string    `bson:"rationale,omitempty" json:"rationale,omitempty"`
	Observation string    `bson:"observation,omitempty" json:"observation,omitempty"`
	Input       string    `bson:"input,omitempty" json:"input,omitempty"`
	Output      string    `bson:"output,omitempty" json:"output,omitempty"`
	StartTime   time.Time `bson:"start_time" json:"startTime"`
	EndTime     time.Time `bson:"end_time,omitempty" json:"endTime,omitempty"`
	Duration    int64     `bson:"duration,omitempty" json:"duration,omitempty"` // in milliseconds
}

type ErrorInfo struct {
	Type       string `bson:"type" json:"type"`
	Message    string `bson:"message" json:"message"`
	Source     string `bson:"source" json:"source"` // Lambda function name
	StackTrace string `bson:"stack_trace,omitempty" json:"stackTrace,omitempty"`
}

// SSE Event types
type SSEEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type ThinkingEvent struct {
	Status string `json:"status"`
}

type AgentStepEvent struct {
	StepIndex   int    `json:"stepIndex"`
	AgentName   string `json:"agentName"`
	AgentID     string `json:"agentId,omitempty"`
	Type        string `json:"type,omitempty"`
	Action      string `json:"action,omitempty"`
	Status      string `json:"status"`
	Rationale   string `json:"rationale,omitempty"`
	Observation string `json:"observation,omitempty"`
	Input       string `json:"input,omitempty"`
	Output      string `json:"output,omitempty"`
	Duration    int64  `json:"duration,omitempty"`
}

type ContentEvent struct {
	Chunk string `json:"chunk"`
}

type TraceEvent struct {
	TraceID    string      `json:"traceId"`
	AgentSteps []AgentStep `json:"agentSteps"`
}

type ErrorEvent struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Source     string `json:"source,omitempty"`
	StackTrace string `json:"stackTrace,omitempty"`
}

type DoneEvent struct {
	MessageID string `json:"messageId"`
}
