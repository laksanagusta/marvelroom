package llm

import (
	"context"
	"fmt"
	"log"

	"sandbox/internal/domain/service"
)

// FallbackLLMService implements LLMService interface with multiple model fallbacks
type FallbackLLMService struct {
	geminiAPIKey  string
	openAIAPIKey  string
	openAIService service.LLMService
}

// NewFallbackLLMService creates a new fallback LLM service with Gemini models and OpenAI GPT-4o-mini as last resort
func NewFallbackLLMService(geminiAPIKey, openAIAPIKey string) (service.LLMService, error) {
	if geminiAPIKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}

	// Initialize OpenAI service if API key is provided
	var openAIService service.LLMService
	var err error
	if openAIAPIKey != "" {
		openAIService, err = NewOpenAIService(openAIAPIKey)
		if err != nil {
			log.Printf("[FallbackLLM] Warning: Failed to initialize OpenAI service: %v", err)
			// Don't fail, just continue without OpenAI fallback
		}
	}

	return &FallbackLLMService{
		geminiAPIKey:  geminiAPIKey,
		openAIAPIKey:  openAIAPIKey,
		openAIService: openAIService,
	}, nil
}

// CheckDocument performs document checking with fallback between models
func (f *FallbackLLMService) CheckDocument(ctx context.Context, req *service.DocumentCheckRequest) (*service.DocumentCheckResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("document check request is required")
	}

	log.Printf("[FallbackLLM] Starting CheckDocument with %d documents", len(req.Documents))

	var lastError error

	// Try each Gemini model in order
	for _, model := range GeminiModels {
		log.Printf("[FallbackLLM] Trying Gemini model: %s", model)

		resp, err := f.tryGeminiModel(ctx, req, model)
		if err != nil {
			log.Printf("[FallbackLLM] Gemini model %s failed: %v", model, err)
			lastError = err
			continue // Try next model
		}

		log.Printf("[FallbackLLM] Successfully used Gemini model: %s", model)
		resp.Model = model
		return resp, nil
	}

	// If all Gemini models failed, try OpenAI GPT-4o-mini as last resort
	if f.openAIService != nil {
		log.Printf("[FallbackLLM] All Gemini models failed, falling back to GPT-4o-mini")

		resp, err := f.openAIService.CheckDocument(ctx, req)
		if err != nil {
			log.Printf("[FallbackLLM] GPT-4o-mini also failed: %v", err)
			lastError = fmt.Errorf("all LLM models failed (last error from GPT-4o-mini: %w)", err)
		} else {
			log.Printf("[FallbackLLM] Successfully used GPT-4o-mini as fallback")
			return resp, nil
		}
	} else {
		log.Printf("[FallbackLLM] OpenAI service not available, cannot fallback to GPT-4o-mini")
	}

	return nil, fmt.Errorf("all LLM models failed: %w", lastError)
}

// tryGeminiModel attempts to use a specific Gemini model
func (f *FallbackLLMService) tryGeminiModel(ctx context.Context, req *service.DocumentCheckRequest, model string) (*service.DocumentCheckResponse, error) {
	// Create a GeminiService with the specific model
	geminiService := &GeminiService{
		apiKey:     f.geminiAPIKey,
		httpClient: LongTimeoutHTTPClient(),
		model:      model,
	}

	return geminiService.CheckDocument(ctx, req)
}
