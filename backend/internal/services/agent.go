package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/ui-agentbedrock/backend/internal/models"
)

type AgentService struct {
	client       *bedrockagentruntime.Client
	awsConfig    aws.Config
	agentID      string
	agentAliasID string
	agentName    string // Display name for the main agent
}

func NewAgentService(agentID, agentAliasID, agentName, region string) (*AgentService, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	client := bedrockagentruntime.NewFromConfig(cfg)

	// Default agent name if not provided
	if agentName == "" {
		agentName = "Main Agent"
	}

	return &AgentService{
		client:       client,
		awsConfig:    cfg,
		agentID:      agentID,
		agentAliasID: agentAliasID,
		agentName:    agentName,
	}, nil
}

// GetAWSConfig returns the AWS configuration for reuse by other services
func (s *AgentService) GetAWSConfig() aws.Config {
	return s.awsConfig
}

type StreamCallback func(event models.SSEEvent) error

func (s *AgentService) InvokeAgentStream(ctx context.Context, sessionID, message string, callback StreamCallback) (*models.Trace, string, error) {
	trace := &models.Trace{
		TraceID:    fmt.Sprintf("trace-%d", time.Now().UnixNano()),
		AgentSteps: []models.AgentStep{},
	}

	var fullContent string
	stepIndex := 0

	// Send thinking event
	callback(models.SSEEvent{
		Event: "thinking",
		Data:  models.ThinkingEvent{Status: "thinking"},
	})

	input := &bedrockagentruntime.InvokeAgentInput{
		AgentId:      aws.String(s.agentID),
		AgentAliasId: aws.String(s.agentAliasID),
		SessionId:    aws.String(sessionID),
		InputText:    aws.String(message),
		EnableTrace:  aws.Bool(true),
		EndSession:   aws.Bool(false),
	}

	output, err := s.client.InvokeAgent(ctx, input)
	if err != nil {
		errorEvent := models.ErrorEvent{
			Type:    "InvokeAgentError",
			Message: err.Error(),
			Source:  "AgentBedrock",
		}
		callback(models.SSEEvent{Event: "error", Data: errorEvent})
		trace.Error = &models.ErrorInfo{
			Type:    errorEvent.Type,
			Message: errorEvent.Message,
			Source:  errorEvent.Source,
		}
		return trace, "", err
	}

	stream := output.GetStream()
	defer stream.Close()

	for event := range stream.Events() {
		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:
			chunk := string(v.Value.Bytes)
			fullContent += chunk
			callback(models.SSEEvent{
				Event: "content",
				Data:  models.ContentEvent{Chunk: chunk},
			})

		case *types.ResponseStreamMemberTrace:
			if v.Value.Trace != nil {
				stepIndex++
				startTime := time.Now()

				// Get agent name from trace - use CollaboratorName if available
				agentName := s.agentName
				if v.Value.CollaboratorName != nil && *v.Value.CollaboratorName != "" {
					agentName = *v.Value.CollaboratorName
				}

				step := s.parseTraceToStep(stepIndex, v.Value.Trace, startTime, agentName)
				if step.Action != "" { // Only add non-empty steps
					trace.AgentSteps = append(trace.AgentSteps, step)

					callback(models.SSEEvent{
						Event: "agent_step",
						Data: models.AgentStepEvent{
							StepIndex:   step.StepIndex,
							AgentName:   step.AgentName,
							AgentID:     step.AgentID,
							Type:        step.Type,
							Action:      step.Action,
							Status:      step.Status,
							Rationale:   step.Rationale,
							Observation: step.Observation,
							Input:       step.Input,
							Output:      step.Output,
							Duration:    step.Duration,
						},
					})
				}
			}
		}
	}

	// Mark all remaining steps as success
	for i := range trace.AgentSteps {
		if trace.AgentSteps[i].Status == "running" {
			trace.AgentSteps[i].Status = "success"
			trace.AgentSteps[i].EndTime = time.Now()
		}
	}

	// Send final agent step status
	for _, step := range trace.AgentSteps {
		callback(models.SSEEvent{
			Event: "agent_step",
			Data: models.AgentStepEvent{
				StepIndex: step.StepIndex,
				AgentName: step.AgentName,
				Action:    step.Action,
				Status:    step.Status,
			},
		})
	}

	// Send trace event
	callback(models.SSEEvent{
		Event: "trace",
		Data: models.TraceEvent{
			TraceID:    trace.TraceID,
			AgentSteps: trace.AgentSteps,
		},
	})

	return trace, fullContent, nil
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// parseTraceToStep extracts detailed information from AWS Bedrock trace
func (s *AgentService) parseTraceToStep(stepIndex int, trace types.Trace, startTime time.Time, defaultAgentName string) models.AgentStep {
	step := models.AgentStep{
		StepIndex: stepIndex,
		AgentName: defaultAgentName, // Use passed agent name (may come from trace)
		AgentID:   s.agentID,
		Type:      "orchestration",
		Action:    "",
		Status:    "success",
		StartTime: startTime,
		EndTime:   time.Now(),
	}
	step.Duration = step.EndTime.Sub(step.StartTime).Milliseconds()

	// Use type switch to handle the Trace interface (outer level)
	switch t := trace.(type) {
	case *types.TraceMemberPreProcessingTrace:
		step.Type = "pre_processing"
		// Keep main agent name but add action context
		step.Action = "Analyzing input"
		// PreProcessingTrace is also an interface - parse inner value
		s.parsePreProcessingTrace(&step, t.Value)

	case *types.TraceMemberOrchestrationTrace:
		step.Type = "orchestration"
		// Keep main agent name - will be updated if it's a collaborator call
		// OrchestrationTrace is also an interface - parse inner value
		s.parseOrchestrationTrace(&step, t.Value)

	case *types.TraceMemberPostProcessingTrace:
		step.Type = "post_processing"
		step.Action = "Formatting response"
		// PostProcessingTrace is also an interface - parse inner value
		s.parsePostProcessingTrace(&step, t.Value)

	case *types.TraceMemberFailureTrace:
		step.Type = "error"
		step.Status = "error"
		if t.Value.FailureReason != nil {
			step.Action = *t.Value.FailureReason
		}

	case *types.TraceMemberGuardrailTrace:
		step.Type = "guardrail"
		step.Action = "Content check"
	}

	return step
}

func (s *AgentService) parsePreProcessingTrace(step *models.AgentStep, trace types.PreProcessingTrace) {
	switch t := trace.(type) {
	case *types.PreProcessingTraceMemberModelInvocationInput:
		if t.Value.Text != nil {
			step.Input = truncateString(*t.Value.Text, 500)
		}
	case *types.PreProcessingTraceMemberModelInvocationOutput:
		if t.Value.ParsedResponse != nil && t.Value.ParsedResponse.Rationale != nil {
			step.Rationale = *t.Value.ParsedResponse.Rationale
		}
	}
}

func (s *AgentService) parseOrchestrationTrace(step *models.AgentStep, trace types.OrchestrationTrace) {
	switch t := trace.(type) {
	case *types.OrchestrationTraceMemberModelInvocationInput:
		if t.Value.Text != nil {
			step.Input = truncateString(*t.Value.Text, 500)
		}
		step.Action = "Processing request"

	case *types.OrchestrationTraceMemberRationale:
		if t.Value.Text != nil {
			// Don't truncate rationale - show full thinking process
			step.Rationale = *t.Value.Text
		}
		step.Action = "Thinking"

	case *types.OrchestrationTraceMemberInvocationInput:
		inv := t.Value

		// Handle Agent Collaborator (Team Agents)
		if inv.AgentCollaboratorInvocationInput != nil {
			collab := inv.AgentCollaboratorInvocationInput
			step.Type = "collaborator"
			step.Status = "running"
			if collab.AgentCollaboratorName != nil {
				step.AgentName = *collab.AgentCollaboratorName
			}
			step.Action = "Calling"
			if collab.Input != nil && collab.Input.Text != nil {
				// Don't truncate collaborator input - show full conversation
				step.Input = *collab.Input.Text
			}
		}

		// Handle Action Group
		if inv.ActionGroupInvocationInput != nil {
			ag := inv.ActionGroupInvocationInput
			step.Type = "action"
			if ag.ActionGroupName != nil {
				step.AgentName = *ag.ActionGroupName
			}
			if ag.Function != nil {
				step.Action = fmt.Sprintf("Function: %s", *ag.Function)
			} else if ag.ApiPath != nil {
				step.Action = fmt.Sprintf("API: %s", *ag.ApiPath)
			} else {
				step.Action = "Executing"
			}
		}

		// Handle Knowledge Base
		if inv.KnowledgeBaseLookupInput != nil {
			kb := inv.KnowledgeBaseLookupInput
			step.Type = "knowledge_base"
			step.AgentName = "Knowledge Base"
			step.Action = "Searching"
			if kb.Text != nil {
				step.Input = *kb.Text
			}
			if kb.KnowledgeBaseId != nil {
				step.AgentID = *kb.KnowledgeBaseId
			}
		}

	case *types.OrchestrationTraceMemberObservation:
		obs := t.Value

		// Handle Agent Collaborator Output (Team Agents)
		if obs.AgentCollaboratorInvocationOutput != nil {
			collab := obs.AgentCollaboratorInvocationOutput
			step.Type = "collaborator"
			if collab.AgentCollaboratorName != nil {
				step.AgentName = *collab.AgentCollaboratorName
			}
			step.Action = "Response"
			if collab.Output != nil && collab.Output.Text != nil {
				// Don't truncate collaborator output - show full conversation
				step.Output = *collab.Output.Text
			}
		}

		// Handle Knowledge Base Output
		if obs.KnowledgeBaseLookupOutput != nil {
			step.Type = "knowledge_base"
			step.AgentName = "Knowledge Base"
			step.Action = "Results"
			if obs.KnowledgeBaseLookupOutput.RetrievedReferences != nil {
				step.Observation = fmt.Sprintf("Found %d references", len(obs.KnowledgeBaseLookupOutput.RetrievedReferences))
			}
		}

		// Handle Action Group Output
		if obs.ActionGroupInvocationOutput != nil && obs.ActionGroupInvocationOutput.Text != nil {
			step.Type = "action"
			step.Action = "Completed"
			step.Output = truncateString(*obs.ActionGroupInvocationOutput.Text, 500)
		}

		// Handle Final Response
		if obs.FinalResponse != nil && obs.FinalResponse.Text != nil {
			step.Action = "Final response"
			step.Output = truncateString(*obs.FinalResponse.Text, 200)
		}

	case *types.OrchestrationTraceMemberModelInvocationOutput:
		step.Action = "Response generated"
	}
}

func (s *AgentService) parsePostProcessingTrace(step *models.AgentStep, trace types.PostProcessingTrace) {
	switch t := trace.(type) {
	case *types.PostProcessingTraceMemberModelInvocationOutput:
		if t.Value.ParsedResponse != nil && t.Value.ParsedResponse.Text != nil {
			step.Output = truncateString(*t.Value.ParsedResponse.Text, 200)
		}
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
