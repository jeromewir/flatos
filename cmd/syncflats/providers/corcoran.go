package providers

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/jeromewir/flatos/internal/config"
	"resty.dev/v3"
)

type CorcoranResponse struct {
	Items      []CorcoranListing `json:"items"`
	PageSize   int               `json:"pageSize"`
	Page       int               `json:"page"`
	TotalPages int               `json:"totalPages"`
	TotalItems int               `json:"totalItems"`
}

type CorcoranListing struct {
	ListingID          string              `json:"listingId"`
	PropertyID         int                 `json:"propertyId"`
	TotalBathrooms     float32             `json:"totalBathrooms"`
	TotalBedrooms      *float32            `json:"totalBedrooms"`
	MaximumLeaseRate   *float64            `json:"maximumLeaseRate"`
	MinimumLeaseRate   *float64            `json:"minimumLeaseRate"`
	Ads                interface{}         `json:"ads"`
	NeighborhoodID     *int                `json:"neighborhoodId"`
	IsHealthCare       bool                `json:"isHealthCare"`
	ClosedRentedDate   string              `json:"closedRentedDate"`
	Location           CorcoranLocation    `json:"location"`
	SourceID           int                 `json:"sourceId"`
	SourceKey          *string             `json:"sourceKey"`
	NewDevelopmentID   *int                `json:"newDevelopmentId"`
	UnitType           *string             `json:"unitType"`
	Agents             []interface{}       `json:"agents"`
	Floorplans         interface{}         `json:"floorplans"`
	IsBuilding         bool                `json:"isBuilding"`
	SquareFootage      float64             `json:"squareFootage"`
	IsNotClickable     bool                `json:"isNotClickable"`
	ID                 int                 `json:"id"`
	AdID               *int                `json:"adID"`
	Source             int                 `json:"source"`
	ListingStatus      string              `json:"listingStatus"`
	ListingStyle       string              `json:"listingStyle"`
	TransactionType    string              `json:"transactionType"`
	RegionID           int                 `json:"regionId"`
	ListingType        string              `json:"listingType"`
	BuildingStyle      *string             `json:"buildingStyle"`
	BuildingType       string              `json:"buildingType"`
	PropertyType       string              `json:"propertyType"`
	Ownership          string              `json:"ownership"`
	Price              float64             `json:"price"`
	Address1           string              `json:"address1"`
	Address2           string              `json:"address2"`
	State              string              `json:"state"`
	ZipCode            string              `json:"zipCode"`
	StreetName         string              `json:"streetName"`
	BoroughName        string              `json:"boroughName"`
	NeighborhoodName   string              `json:"neighborhoodName"`
	City               string              `json:"city"`
	Bedrooms           *float32            `json:"bedrooms"`
	Bathrooms          float32             `json:"bathrooms"`
	HalfBaths          *float32            `json:"halfBaths"`
	IsNew              bool                `json:"isNew"`
	IsVow              bool                `json:"isVow"`
	IsReducedPrice     bool                `json:"isReducedPrice"`
	IsIdx              bool                `json:"isIdx"`
	IsValue            bool                `json:"isValue"`
	IsFeatured         bool                `json:"isFeatured"`
	IsAddressHidden    bool                `json:"isAddressHidden"`
	IsPriceHidden      bool                `json:"isPriceHidden"`
	IsFranchise        bool                `json:"isFranchise"`
	IsExclusive        bool                `json:"isExclusive"`
	HasVirtualTour     bool                `json:"hasVirtualTour"`
	IsOpenHouse        bool                `json:"isOpenHouse"`
	IsVirtualOpenHouse bool                `json:"isVirtualOpenHouse"`
	AdvertiseNoFee     bool                `json:"advertiseNoFee"`
	MediaUrl           *string             `json:"mediaUrl"`
	Media              []interface{}       `json:"media"`
	WatermarkLogoUrl   *string             `json:"watermarkLogoUrl"`
	OpenHouses         []interface{}       `json:"openHouses"`
	VirtualOpenHouses  []interface{}       `json:"virtualOpenHouses"`
	SquareFeetInterior *float64            `json:"squareFeetInterior"`
	InteriorSquare     CorcoranInterior    `json:"interiorSquare"`
	ListingAttribution CorcoranAttribution `json:"listingAttribution"`
	ListingMlsNumber   *string             `json:"listingMlsNumber"`
	IsTrending         bool                `json:"isTrending"`
	ListedDate         string              `json:"listedDate"`
	DateAvailable      string              `json:"dateAvailable"`
	IsMLS              *bool               `json:"isMLS"`
	IsNewDevelopment   bool                `json:"isNewDevelopment"`
	ListingAgency      string              `json:"listingAgency"`
	AvailableDate      *string             `json:"availableDate"`
	Country            string              `json:"country"`
}

type CorcoranLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CorcoranInterior struct {
	Value         float64 `json:"value"`
	UnitOfMeasure string  `json:"unitOfMeasure"`
	SystemType    int     `json:"systemType"`
}

type CorcoranAttribution struct {
	CompanyWebID       *string `json:"companyWebId"`
	CompanyLogoUrl     *string `json:"companyLogoUrl"`
	CompanyName        string  `json:"companyName"`
	ExternalCompanyUrl *string `json:"externalCompanyUrl"`
	IsClickable        bool    `json:"isClickable"`
}

type Corcoran struct {
	config *config.Config
}

func NewCorcoran(config *config.Config) *Corcoran {
	return &Corcoran{
		config: config,
	}
}

func (c *Corcoran) extractIDFromURL(url string) (string, error) {
	// Example: https://www.corcoran.com/building/upper-west-side/4806
	re := `\/(\d+)$`
	matches := regexp.MustCompile(re).FindStringSubmatch(url)

	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract ID from url: %s", url)
	}

	return matches[1], nil
}

func (c *Corcoran) GetFlats(url string) ([]Flat, error) {
	restyClient := resty.New()
	var response CorcoranResponse

	id, err := c.extractIDFromURL(url)

	if err != nil {
		return nil, err
	}

	resp, err := restyClient.
		NewRequest().
		SetResult(&response).
		SetHeader("be-api-key", c.config.CorcoranAPIKey).
		Get(
			fmt.Sprintf("https://backendapi.corcoranlabs.com/api/properties/%s/listings?page=1&pageSize=30&tabName=rent&currency=USD&measure=imperial&languageCode=en-US", id),
		)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get flats: %s", resp.Status())
	}

	flats := make([]Flat, len(response.Items))
	for i, item := range response.Items {
		createdAt, err := time.Parse(time.RFC3339, item.ListedDate)

		if err != nil {
			createdAt = time.Time{}
		}

		flats[i] = Flat{
			ID:           fmt.Sprintf("flt_%s", uuid.New().String()),
			ExternalID:   fmt.Sprintf("%d", item.ID),
			Name:         item.Address1,
			Address:      item.Address1,
			City:         item.City,
			Price:        int64(item.Price),
			ZipCode:      item.ZipCode,
			Size:         fmt.Sprintf("%f", item.SquareFootage),
			Availability: item.ListingStatus,
			Status:       item.ListingStyle,
			CreatedAt:    createdAt,
			UpdatedAt:    time.Now(),
		}
	}

	return flats, nil
}

// The line below is a compile-time assertion to ensure that Corcoran implements the Provider interface.
var _ Provider = (*Corcoran)(nil)
