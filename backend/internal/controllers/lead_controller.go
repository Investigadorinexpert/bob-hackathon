package controllers

import (
	"bob-hackathon/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LeadController struct {
	sessionService *services.SessionService
	faqService     *services.FAQService
	bobAPIService  *services.BOBAPIService
}

func NewLeadController() *LeadController {
	return &LeadController{
		sessionService: services.GetSessionService(),
		faqService:     services.GetFAQService(),
		bobAPIService:  services.GetBOBAPIService(),
	}
}

func (l *LeadController) GetAllLeads(ctx *gin.Context) {
	category := ctx.Query("category")
	channel := ctx.Query("channel")

	leads := l.sessionService.GetAllLeads(category, channel)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(leads),
		"leads":   leads,
	})
}

func (l *LeadController) GetLead(ctx *gin.Context) {
	sessionID := ctx.Param("sessionId")

	lead := l.sessionService.GetLead(sessionID)
	if lead == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Lead no encontrado",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"lead":    lead,
	})
}

func (l *LeadController) GetLeadsStats(ctx *gin.Context) {
	stats := l.sessionService.GetLeadsStats()

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats":   stats,
	})
}

func (l *LeadController) GetFAQs(ctx *gin.Context) {
	search := ctx.Query("search")
	categoria := ctx.Query("categoria")
	empresa := ctx.Query("empresa")

	faqs := l.faqService.SearchFAQs(search, categoria, empresa)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(faqs),
		"faqs":    faqs,
	})
}

func (l *LeadController) GetVehicles(ctx *gin.Context) {
	marca := ctx.Query("marca")
	modelo := ctx.Query("modelo")
	tipoSubasta := ctx.Query("tipo_subasta")

	precioMinStr := ctx.Query("precio_min")
	precioMaxStr := ctx.Query("precio_max")
	limitStr := ctx.DefaultQuery("limit", "10")

	precioMin := 0.0
	precioMax := 0.0
	limit := 10

	if precioMinStr != "" {
		if val, err := strconv.ParseFloat(precioMinStr, 64); err == nil {
			precioMin = val
		}
	}

	if precioMaxStr != "" {
		if val, err := strconv.ParseFloat(precioMaxStr, 64); err == nil {
			precioMax = val
		}
	}

	if val, err := strconv.Atoi(limitStr); err == nil {
		limit = val
	}

	vehicles, err := l.bobAPIService.SearchVehicles(marca, modelo, precioMin, precioMax, tipoSubasta, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error al obtener vehículos: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"count":    len(vehicles),
		"vehicles": vehicles,
	})
}

func (l *LeadController) GetVehicleByID(ctx *gin.Context) {
	id := ctx.Param("id")

	vehicle, err := l.bobAPIService.GetVehicleByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Vehículo no encontrado",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"vehicle": vehicle,
	})
}
