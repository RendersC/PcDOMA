package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PCAvailability struct {
	Available    bool    `json:"available"`
	PricePerHour float64 `json:"price_per_hour"`
	PricePerDay  float64 `json:"price_per_day"`
}

type SetupDetails struct {
	ID                string   `json:"id"`
	PCID              string   `json:"pc_id"`
	PeripheralIDs     []string `json:"peripheral_ids"`
	TotalPricePerHour float64  `json:"total_price_per_hour"`
	TotalPricePerDay  float64  `json:"total_price_per_day"`
}

type CatalogClient struct {
	baseURL string
	client  *http.Client
}

func NewCatalogClient(baseURL string) *CatalogClient {
	return &CatalogClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *CatalogClient) GetPCAvailability(ctx context.Context, pcID string) (*PCAvailability, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/internal/pcs/%s/availability", c.baseURL, pcID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("catalog service returned %d", resp.StatusCode)
	}

	var avail PCAvailability
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &avail); err != nil {
		return nil, err
	}
	return &avail, nil
}

func (c *CatalogClient) GetSetupDetails(ctx context.Context, setupID string) (*SetupDetails, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/catalog/setups/%s", c.baseURL, setupID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("catalog service returned %d", resp.StatusCode)
	}
	var setup SetupDetails
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &setup); err != nil {
		return nil, err
	}
	return &setup, nil
}

func (c *CatalogClient) UpdatePCStatus(ctx context.Context, pcID, status string) error {
	body, _ := json.Marshal(map[string]string{"status": status})
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch,
		fmt.Sprintf("%s/internal/pcs/%s/status", c.baseURL, pcID),
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("catalog service returned %d", resp.StatusCode)
	}
	return nil
}
