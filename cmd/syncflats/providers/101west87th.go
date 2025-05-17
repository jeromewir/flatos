package providers

import (
	"github.com/jeromewir/flatos/internal/config"
)

type OneOOneWestEightySeven struct {
	ESDataProvider
}

func NewOneOOneWestEightySeven(config *config.Config) *OneOOneWestEightySeven {
	return &OneOOneWestEightySeven{
		ESDataProvider: ESDataProvider{
			config:  config,
			address: "101 West 87th Street",
			zipCode: "10024",
		},
	}
}

func (c *OneOOneWestEightySeven) GetFlats(url string) ([]Flat, error) {
	flats, err := c.ESDataProvider.GetFlats(url)

	if err != nil {
		return nil, err
	}

	return flats, nil
}

// The line below is a compile-time assertion to ensure that OneOOneWestEightySeven implements the Provider interface.
var _ Provider = (*OneOOneWestEightySeven)(nil)
