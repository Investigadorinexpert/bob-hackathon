package controllers

import (
	"bob-hackathon/internal/models"
	"bob-hackathon/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	geminiService  *services.GeminiService
	sessionService *services.SessionService
}

func NewChatController() *ChatController {
	return &ChatController{
		geminiService:  services.GetGeminiService(),
		sessionService: services.GetSessionService(),
	}
}

func (c *ChatController) SendMessage(ctx *gin.Context) {
	var req models.ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Obtener o crear sesión
	session := c.sessionService.GetOrCreateSession(req.SessionID, req.Channel)

	// Agregar mensaje del usuario
	c.sessionService.AddMessage(session.SessionID, "user", req.Message)

	// Procesar con Gemini
	reply, err := c.geminiService.ProcessMessage(session.SessionID, req.Message)
	if err != nil {
		log.Printf("❌ Error al procesar mensaje: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error al procesar mensaje: " + err.Error(),
		})
		return
	}

	// Agregar respuesta del asistente
	c.sessionService.AddMessage(session.SessionID, "assistant", reply)

	// Calcular score
	scoreResponse, err := c.geminiService.CalculateScore(session.SessionID)
	if err != nil {
		log.Printf("Error al calcular score: %v", err)
		scoreResponse = &models.ScoreResponse{
			Score:    0,
			Category: "cold",
		}
	}

	// Actualizar score en sesión
	c.sessionService.UpdateScore(session.SessionID, scoreResponse.Score, scoreResponse.Category)

	// Crear o actualizar lead
	lead := &models.Lead{
		SessionID:    session.SessionID,
		Channel:      session.Channel,
		Score:        scoreResponse.Score,
		Category:     scoreResponse.Category,
		Urgency:      scoreResponse.Urgency,
		Budget:       scoreResponse.Budget,
		BusinessType: scoreResponse.BusinessType,
		Reasons:      scoreResponse.Reasons,
		LastMessage:  req.Message,
	}
	c.sessionService.CreateOrUpdateLead(lead)

	// Responder
	response := models.ChatResponse{
		Success:   true,
		SessionID: session.SessionID,
		Reply:     reply,
		LeadScore: scoreResponse.Score,
		Category:  scoreResponse.Category,
		Timestamp: time.Now(),
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ChatController) GetScore(ctx *gin.Context) {
	var req models.ScoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Datos inválidos: " + err.Error(),
		})
		return
	}

	scoreResponse, err := c.geminiService.CalculateScore(req.SessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error al calcular score: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, scoreResponse)
}

func (c *ChatController) GetHistory(ctx *gin.Context) {
	sessionID := ctx.Param("sessionId")

	session := c.sessionService.GetSession(sessionID)
	if session == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Sesión no encontrada",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"session":  session,
		"messages": session.Messages,
	})
}

func (c *ChatController) DeleteSession(ctx *gin.Context) {
	_ = ctx.Param("sessionId") // sessionID no usado aún

	// Por ahora solo retornamos éxito
	// Se puede implementar eliminación si es necesario
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sesión eliminada (funcionalidad pendiente)",
	})
}
