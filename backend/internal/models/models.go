package models

import "time"

// Session representa una sesión de conversación
type Session struct {
	SessionID    string              `json:"sessionId"`
	Channel      string              `json:"channel"`
	Messages     []Message           `json:"messages"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
	LeadScore    int                 `json:"leadScore"`
	Category     string              `json:"category"`
	Metadata     map[string]string   `json:"metadata,omitempty"`
}

// Message representa un mensaje en la conversación
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Lead representa un lead generado
type Lead struct {
	SessionID    string              `json:"sessionId"`
	Channel      string              `json:"channel"`
	Score        int                 `json:"score"`
	Category     string              `json:"category"`
	Urgency      string              `json:"urgency,omitempty"`
	Budget       string              `json:"budget,omitempty"`
	BusinessType string              `json:"businessType,omitempty"`
	Reasons      []string            `json:"reasons,omitempty"`
	LastMessage  string              `json:"lastMessage"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
	Metadata     map[string]string   `json:"metadata,omitempty"`
}

// FAQ representa una pregunta frecuente
type FAQ struct {
	Categoria string `json:"categoria"`
	Empresa   string `json:"empresa"`
	Pregunta  string `json:"pregunta"`
	Respuesta string `json:"respuesta"`
}

// Vehicle representa un vehículo en subasta
type Vehicle struct {
	ID           string  `json:"id"`
	Marca        string  `json:"marca"`
	Modelo       string  `json:"modelo"`
	Ano          string  `json:"ano"`
	PrecioInicio float64 `json:"precioInicio"`
	TipoSubasta  string  `json:"tipoSubasta"`
	Estado       string  `json:"estado"`
	Imagen       string  `json:"imagen,omitempty"`
}

// ChatRequest representa una solicitud de mensaje
type ChatRequest struct {
	SessionID string `json:"sessionId,omitempty"`
	Message   string `json:"message" binding:"required"`
	Channel   string `json:"channel" binding:"required"`
}

// ChatResponse representa la respuesta del chat
type ChatResponse struct {
	Success   bool      `json:"success"`
	SessionID string    `json:"sessionId"`
	Reply     string    `json:"reply"`
	LeadScore int       `json:"leadScore"`
	Category  string    `json:"category"`
	Timestamp time.Time `json:"timestamp"`
}

// ScoreRequest representa una solicitud de scoring
type ScoreRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
}

// ScoreResponse representa la respuesta de scoring
type ScoreResponse struct {
	Success      bool     `json:"success"`
	Score        int      `json:"score"`
	Category     string   `json:"category"`
	Reasons      []string `json:"reasons"`
	Urgency      string   `json:"urgency"`
	Budget       string   `json:"budget"`
	BusinessType string   `json:"businessType"`
}

// LeadStats representa estadísticas de leads
type LeadStats struct {
	Total      int     `json:"total"`
	Hot        int     `json:"hot"`
	Warm       int     `json:"warm"`
	Cold       int     `json:"cold"`
	AvgScore   float64 `json:"avgScore"`
	ByChannel  map[string]int `json:"byChannel"`
}

// HealthResponse representa la respuesta del health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}
