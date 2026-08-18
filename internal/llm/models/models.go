package models

import "maps"

type (
	ModelID       string
	ModelProvider string
)

type Model struct {
	ID                  ModelID       `json:"id"`
	Name                string        `json:"name"`
	Provider            ModelProvider `json:"provider"`
	APIModel            string        `json:"api_model"`
	CostPer1MIn         float64       `json:"cost_per_1m_in"`
	CostPer1MOut        float64       `json:"cost_per_1m_out"`
	CostPer1MInCached   float64       `json:"cost_per_1m_in_cached"`
	CostPer1MOutCached  float64       `json:"cost_per_1m_out_cached"`
	ContextWindow       int64         `json:"context_window"`
	DefaultMaxTokens    int64         `json:"default_max_tokens"`
	CanReason           bool          `json:"can_reason"`
	SupportsAttachments bool          `json:"supports_attachments"`
}

// Model IDs
const ( // GEMINI
	// Bedrock
	BedrockClaude37Sonnet ModelID = "bedrock.claude-3.7-sonnet"
)

const (
	ProviderBedrock ModelProvider = "bedrock"
	// ForTests
	ProviderMock ModelProvider = "__mock"
)

// Providers in order of popularity
var ProviderPopularity = map[ModelProvider]int{
	ProviderCopilot:    1,
	ProviderAnthropic:  2,
	ProviderOpenAI:     3,
	ProviderGemini:     4,
	ProviderGROQ:       5,
	ProviderOpenRouter: 6,
	ProviderBedrock:    7,
	ProviderAzure:      8,
	ProviderVertexAI:   9,
}

var SupportedModels = map[ModelID]Model{
	//
	// // GEMINI
	// GEMINI25: {
	// 	ID:                 GEMINI25,
	// 	Name:               "Gemini 2.5 Pro",
	// 	Provider:           ProviderGemini,
	// 	APIModel:           "gemini-2.5-pro-exp-03-25",
	// 	CostPer1MIn:        0,
	// 	CostPer1MInCached:  0,
	// 	CostPer1MOutCached: 0,
	// 	CostPer1MOut:       0,
	// },
	//
	// GRMINI20Flash: {
	// 	ID:                 GRMINI20Flash,
	// 	Name:               "Gemini 2.0 Flash",
	// 	Provider:           ProviderGemini,
	// 	APIModel:           "gemini-2.0-flash",
	// 	CostPer1MIn:        0.1,
	// 	CostPer1MInCached:  0,
	// 	CostPer1MOutCached: 0.025,
	// 	CostPer1MOut:       0.4,
	// },
	//
	// // Bedrock
	BedrockClaude37Sonnet: {
		ID:                 BedrockClaude37Sonnet,
		Name:               "Bedrock: Claude 3.7 Sonnet",
		Provider:           ProviderBedrock,
		APIModel:           "anthropic.claude-3-7-sonnet-20250219-v1:0",
		CostPer1MIn:        3.0,
		CostPer1MInCached:  3.75,
		CostPer1MOutCached: 0.30,
		CostPer1MOut:       15.0,
	},
}

func init() {
	maps.Copy(SupportedModels, AnthropicModels)
	maps.Copy(SupportedModels, OpenAIModels)
	maps.Copy(SupportedModels, GeminiModels)
	maps.Copy(SupportedModels, GroqModels)
	maps.Copy(SupportedModels, AzureModels)
	maps.Copy(SupportedModels, OpenRouterModels)
	maps.Copy(SupportedModels, XAIModels)
	maps.Copy(SupportedModels, VertexAIGeminiModels)
	maps.Copy(SupportedModels, CopilotModels)
}

// dynamicModels holds models discovered from custom (OpenAI-compatible)
// providers at runtime, e.g. via the /provider command + /models dialog.
var dynamicModels = make(map[ModelID]Model)

// RegisterDynamicModel adds a model discovered from a custom provider endpoint.
func RegisterDynamicModel(id ModelID, name string, provider ModelProvider) {
	if id == "" {
		return
	}
	dynamicModels[id] = Model{
		ID:               id,
		Name:             name,
		Provider:         provider,
		APIModel:         string(id),
		ContextWindow:    128000,
		DefaultMaxTokens: 4096,
	}
}

// RegisterDynamicModels adds a list of model IDs discovered from a custom
// provider endpoint.
func RegisterDynamicModels(ids []string, provider ModelProvider) {
	for _, id := range ids {
		RegisterDynamicModel(ModelID(id), id, provider)
	}
}

// GetModel looks up a model in the built-in or dynamically registered models.
func GetModel(id ModelID) (Model, bool) {
	if m, ok := SupportedModels[id]; ok {
		return m, true
	}
	m, ok := dynamicModels[id]
	return m, ok
}

// ModelName returns a display name for a model ID, falling back to the ID
// itself when the model is unknown (e.g. a custom provider model).
func ModelName(id ModelID) string {
	if m, ok := GetModel(id); ok {
		return m.Name
	}
	return string(id)
}

// ModelsForProvider returns all models (built-in + dynamic) for a provider.
func ModelsForProvider(provider ModelProvider) []Model {
	var result []Model
	for _, m := range SupportedModels {
		if m.Provider == provider {
			result = append(result, m)
		}
	}
	for _, m := range dynamicModels {
		if m.Provider == provider {
			result = append(result, m)
		}
	}
	return result
}
