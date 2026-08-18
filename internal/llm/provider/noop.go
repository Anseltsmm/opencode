package provider

import (
	"context"
	"errors"

	"github.com/Anseltsmm/azkia/internal/llm/models"
	"github.com/Anseltsmm/azkia/internal/llm/tools"
	"github.com/Anseltsmm/azkia/internal/message"
)

// ErrNoProviderConfigured is returned when no LLM provider has been configured
// yet (no API key / base URL in config or environment).
var ErrNoProviderConfigured = errors.New("no provider configured")

type noopProvider struct{}

func (p *noopProvider) SendMessages(ctx context.Context, messages []message.Message, tools []tools.BaseTool) (*ProviderResponse, error) {
	return nil, ErrNoProviderConfigured
}

func (p *noopProvider) StreamResponse(ctx context.Context, messages []message.Message, tools []tools.BaseTool) <-chan ProviderEvent {
	ch := make(chan ProviderEvent, 1)
	ch <- ProviderEvent{Type: EventError, Error: ErrNoProviderConfigured}
	close(ch)
	return ch
}

func (p *noopProvider) Model() models.Model {
	return models.Model{Name: "no provider"}
}

// NewNoopProvider returns a provider that errors when used. It lets the TUI
// start without a configured provider so the user can be prompted to set one up.
func NewNoopProvider() Provider {
	return &noopProvider{}
}
