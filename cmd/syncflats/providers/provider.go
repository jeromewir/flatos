package providers

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jeromewir/flatos/internal/config"
)

type Flat struct {
	ID           string
	ExternalID   string
	Name         string
	Address      string
	City         string
	Price        int64
	ZipCode      string
	Size         string
	Availability string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Provider interface {
	GetFlats(url string) ([]Flat, error)
}

func Get(config *config.Config, urlString string) (Provider, error) {
	parsedURL, err := url.Parse(urlString)

	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	host := strings.Replace(parsedURL.Host, "www.", "", 1)

	switch host {
	case CorcoranHost:
		return NewCorcoran(config), nil
	case EquityApartmentsHost:
		return NewEquityApartments(config), nil
	case OneOOneWestEightySevenHost:
		return NewOneOOneWestEightySeven(config), nil
	case HLResHost:
		return NewHLRes(config), nil
	}

	return nil, fmt.Errorf("provider not found")
}
