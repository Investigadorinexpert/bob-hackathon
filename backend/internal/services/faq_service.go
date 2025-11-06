package services

import (
	"bob-hackathon/internal/models"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FAQService struct {
	faqs []models.FAQ
	mu   sync.RWMutex
}

var faqServiceInstance *FAQService
var faqServiceOnce sync.Once

func GetFAQService() *FAQService {
	faqServiceOnce.Do(func() {
		faqServiceInstance = &FAQService{
			faqs: []models.FAQ{},
		}
		faqServiceInstance.loadFAQs()
	})
	return faqServiceInstance
}

func (f *FAQService) loadFAQs() {
	file, err := os.Open(filepath.Join("data", "faqs.csv"))
	if err != nil {
		log.Printf("Error al abrir FAQs: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error al leer FAQs: %v", err)
		return
	}

	// Saltar header
	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) >= 4 {
			faq := models.FAQ{
				Categoria: record[0],
				Empresa:   record[1],
				Pregunta:  record[2],
				Respuesta: record[3],
			}
			f.faqs = append(f.faqs, faq)
		}
	}

	log.Printf("%d FAQs cargadas", len(f.faqs))
}

func (f *FAQService) SearchFAQs(query, categoria, empresa string) []models.FAQ {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var results []models.FAQ
	queryLower := strings.ToLower(query)

	for _, faq := range f.faqs {
		// Filtrar por categoría si se especifica
		if categoria != "" && !strings.EqualFold(faq.Categoria, categoria) {
			continue
		}

		// Filtrar por empresa si se especifica
		if empresa != "" && !strings.EqualFold(faq.Empresa, empresa) {
			continue
		}

		// Si hay query, buscar en pregunta y respuesta
		if query != "" {
			preguntaLower := strings.ToLower(faq.Pregunta)
			respuestaLower := strings.ToLower(faq.Respuesta)

			if !strings.Contains(preguntaLower, queryLower) && !strings.Contains(respuestaLower, queryLower) {
				continue
			}
		}

		results = append(results, faq)
	}

	return results
}

func (f *FAQService) GetAllFAQs() []models.FAQ {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.faqs
}

func (f *FAQService) GetFAQsContext() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var context strings.Builder
	context.WriteString("Base de conocimiento (FAQs):\n\n")

	for i, faq := range f.faqs {
		if i >= 10 { // Limitar a 10 FAQs para el contexto
			break
		}
		context.WriteString("P: ")
		context.WriteString(faq.Pregunta)
		context.WriteString("\nR: ")
		context.WriteString(faq.Respuesta)
		context.WriteString("\n\n")
	}

	return context.String()
}
