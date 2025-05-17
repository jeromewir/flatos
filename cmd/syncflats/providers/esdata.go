package providers

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jeromewir/flatos/internal/config"
	"resty.dev/v3"
)

type ESDataListingDetailResponse struct {
	PK                   int      `json:"pk"`
	Photos               []string `json:"photos"`
	Courtesy             string   `json:"courtesy"`
	Unit                 string   `json:"unit"`
	Price                string   `json:"price"`
	Beds                 string   `json:"beds"`
	Baths                string   `json:"baths"`
	SquareFeet           int      `json:"square feet"`
	CommonChargesOrMaint string   `json:"common charges/maintenance"`
	PropertyTaxes        string   `json:"property taxes"`
	CtaText              string   `json:"cta_text"`
	CtaValue             string   `json:"cta_value"`
	Status               string   `json:"status"`
}

// This is a generic implementation as this is used by multiple buildings
type ESDataProvider struct {
	config  *config.Config
	address string
	zipCode string
}

func NewESDataProvider(config *config.Config, address string, zipCode string) *ESDataProvider {
	return &ESDataProvider{
		config:  config,
		address: address,
		zipCode: zipCode,
	}
}

func (c *ESDataProvider) GetFlats(url string) ([]Flat, error) {
	restyClient := resty.New()
	var response []ESDataListingDetailResponse

	resp, err := restyClient.
		NewRequest().
		SetResult(&response).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get flats: %s", resp.Status())
	}

	flats := make([]Flat, len(response))
	for i, item := range response {
		price, _ := strconv.Atoi(regexp.MustCompile("\\$|,").ReplaceAllString(item.Price, ""))

		flats[i] = Flat{
			ID:           fmt.Sprintf("flt_%s", uuid.New().String()),
			ExternalID:   fmt.Sprintf("%d", item.PK),
			Name:         item.Unit,
			Address:      c.address,
			City:         "New York",
			Price:        int64(price),
			ZipCode:      c.zipCode,
			Size:         fmt.Sprintf("%d", item.SquareFeet),
			Availability: item.Status,
			Status:       item.Status,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
	}

	return flats, nil
}

// The line below is a compile-time assertion to ensure that ESDataProvider implements the Provider interface.
var _ Provider = (*ESDataProvider)(nil)
