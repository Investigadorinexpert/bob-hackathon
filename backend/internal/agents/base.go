package agents

import (
	"bob-hackathon/internal/models"
	"context"
)

type Agent interface {
	Process(ctx context.Context, input *AgentInput) (*AgentOutput, error)
	Name() string
}

type AgentInput struct {
	Message        string
	SessionID      string
	Channel        string
	ConversationHistory []models.Message
	LeadData       *models.LeadData
}

type AgentOutput struct {
	Response       string
	ShouldRoute    bool
	RouteTo        string
	ScoringData    *models.ScoringData
	IntentDetected string
	Confidence     float64
}

type IntentType string

const (
	IntentFAQ      IntentType = "faq"
	IntentSubasta  IntentType = "subasta"
	IntentSpam     IntentType = "spam"
	IntentAmbiguo  IntentType = "ambiguo"
	IntentGeneral  IntentType = "general"
)
