package controllers

import (
	"bob-hackathon/internal/services"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	faqService *services.FAQService
}

func NewAdminController(faqService *services.FAQService) *AdminController {
	return &AdminController{
		faqService: faqService,
	}
}

// UploadFAQs maneja la subida de CSV de FAQs
func (a *AdminController) UploadFAQs(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "no se encontro archivo en el request",
		})
		return
	}

	// Validar que sea CSV
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "el archivo debe ser CSV",
		})
		return
	}

	// Validar CSV antes de guardar
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "error al abrir archivo",
		})
		return
	}
	defer src.Close()

	reader := csv.NewReader(src)
	records, err := reader.ReadAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "error al parsear CSV: " + err.Error(),
		})
		return
	}

	// Validar formato (debe tener al menos header + 1 fila)
	if len(records) < 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "CSV debe tener al menos una fila de datos",
		})
		return
	}

	// Validar que tenga 4 columnas
	if len(records[0]) != 4 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "CSV debe tener 4 columnas: categoria,empresa,pregunta,respuesta",
		})
		return
	}

	// Guardar archivo en data/faqs.csv
	destPath := filepath.Join("data", "faqs.csv")

	// Crear backup del archivo anterior
	if _, err := os.Stat(destPath); err == nil {
		backupPath := filepath.Join("data", "faqs.csv.backup")
		os.Rename(destPath, backupPath)
		log.Printf("Backup creado: %s", backupPath)
	}

	// Guardar nuevo archivo
	if err := ctx.SaveUploadedFile(file, destPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "error al guardar archivo",
		})
		return
	}

	// Recargar FAQs en memoria
	services.ReloadFAQs()

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "FAQs actualizadas correctamente",
		"total":     len(records) - 1,
		"timestamp": ctx.GetTime("timestamp"),
	})
}

// GetPrompts devuelve los prompts de todos los agentes
func (a *AdminController) GetPrompts(ctx *gin.Context) {
	prompts := make(map[string]string)

	agents := []string{"orchestrator", "faq", "auction", "scoring"}

	for _, agent := range agents {
		promptPath := filepath.Join("data", "prompts", agent+".txt")
		content, err := os.ReadFile(promptPath)
		if err != nil {
			// Si no existe, devolver mensaje indicativo
			prompts[agent] = fmt.Sprintf("Prompt no configurado (usar default del sistema)")
			continue
		}
		prompts[agent] = string(content)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"prompts": prompts,
	})
}

// UpdatePrompt actualiza el prompt de un agente específico
func (a *AdminController) UpdatePrompt(ctx *gin.Context) {
	agentName := ctx.Param("agent")

	// Validar que sea un agente válido
	validAgents := map[string]bool{
		"orchestrator": true,
		"faq":          true,
		"auction":      true,
		"scoring":      true,
	}

	if !validAgents[agentName] {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "agente invalido. Debe ser: orchestrator, faq, auction, scoring",
		})
		return
	}

	var req struct {
		Prompt string `json:"prompt" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "prompt es requerido",
		})
		return
	}

	// Validar que el prompt no esté vacío
	if strings.TrimSpace(req.Prompt) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "prompt no puede estar vacio",
		})
		return
	}

	// Crear directorio si no existe
	promptsDir := filepath.Join("data", "prompts")
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "error al crear directorio de prompts",
		})
		return
	}

	// Guardar prompt
	promptPath := filepath.Join(promptsDir, agentName+".txt")

	// Crear backup si existe
	if _, err := os.Stat(promptPath); err == nil {
		backupPath := promptPath + ".backup"
		os.Rename(promptPath, backupPath)
		log.Printf("Backup de prompt creado: %s", backupPath)
	}

	if err := os.WriteFile(promptPath, []byte(req.Prompt), 0644); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "error al guardar prompt",
		})
		return
	}

	log.Printf("Prompt de %s actualizado", agentName)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Prompt de %s actualizado correctamente", agentName),
		"agent":   agentName,
	})
}

// DownloadFAQsTemplate descarga un template CSV de FAQs
func (a *AdminController) DownloadFAQsTemplate(ctx *gin.Context) {
	// Template con formato correcto
	template := [][]string{
		{"categoria", "empresa", "pregunta", "respuesta"},
		{"general", "bob", "¿Cómo funciona el proceso de subasta?", "El proceso de subasta en BOB es simple y transparente..."},
		{"vehiculos", "bob", "¿Qué tipos de vehículos ofrecen?", "Ofrecemos una amplia variedad de vehículos..."},
	}

	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename=faqs_template.csv")

	writer := csv.NewWriter(ctx.Writer)
	defer writer.Flush()

	for _, row := range template {
		if err := writer.Write(row); err != nil {
			log.Printf("Error escribiendo CSV template: %v", err)
			return
		}
	}
}

// GetFAQsAsCSV devuelve las FAQs actuales como CSV para descarga
func (a *AdminController) GetFAQsAsCSV(ctx *gin.Context) {
	faqs := a.faqService.GetAllFAQs()

	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename=faqs_current.csv")

	writer := csv.NewWriter(ctx.Writer)
	defer writer.Flush()

	// Header
	writer.Write([]string{"categoria", "empresa", "pregunta", "respuesta"})

	// Datos
	for _, faq := range faqs {
		writer.Write([]string{
			faq.Categoria,
			faq.Empresa,
			faq.Pregunta,
			faq.Respuesta,
		})
	}
}
