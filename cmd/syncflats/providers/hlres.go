package providers

import (
	"fmt"
	"regexp"

	"github.com/jeromewir/flatos/internal/config"
)

type HLRes struct {
	ESDataProvider
}

func NewHLRes(config *config.Config) *HLRes {
	return &HLRes{
		ESDataProvider: ESDataProvider{
			config: config,
			// TODO: find a way to get the address and zip code from the URL
			address: "",
			zipCode: "",
		},
	}
}

func (c *HLRes) extractID(url string) (string, error) {
	r := regexp.MustCompile(`(\d+)\/condo`)

	matches := r.FindStringSubmatch(url)

	if len(matches) < 2 {
		return "", fmt.Errorf("invalid URL format: %s", url)
	}

	id := matches[1]

	if id == "" {
		return "", fmt.Errorf("ID not found in URL: %s", url)
	}

	return id, nil
}

func (c *HLRes) GetFlats(url string) ([]Flat, error) {
	id, err := c.extractID(url)

	if err != nil {
		return nil, err
	}

	flats, err := c.ESDataProvider.GetFlats(fmt.Sprintf("https://www.hlres.com/buildings/ajax/es-data?pk=%s&type=rentals", id))

	if err != nil {
		return nil, err
	}

	return flats, nil
}

// The line below is a compile-time assertion to ensure that HLRes implements the Provider interface.
var _ Provider = (*HLRes)(nil)
