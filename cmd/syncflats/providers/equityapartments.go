package providers

import (
	"github.com/jeromewir/flatos/internal/config"
)

type EquityApartments struct {
	config *config.Config
}

func NewEquityApartments(config *config.Config) *EquityApartments {
	return &EquityApartments{
		config: config,
	}
}

func (e *EquityApartments) GetFlats(url string) ([]Flat, error) {
	// Uses cloudflare

	return []Flat{}, nil
}

var _ Provider = (*EquityApartments)(nil)
