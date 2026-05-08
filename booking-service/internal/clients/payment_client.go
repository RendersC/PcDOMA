package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CreatePaymentResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type PaymentClient struct {
	baseURL string
	client  *http.Client
}

func NewPaymentClient(baseURL string) *PaymentClient {
	return &PaymentClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *PaymentClient) CreatePayment(ctx context.Context, bookingID, userID uuid.UUID, amount float64, currency string) (*CreatePaymentResponse, error) {
	if currency == "" {
		currency = "KZT"
	}
	body, _ := json.Marshal(map[string]interface{}{
		"booking_id": bookingID,
		"user_id":    userID,
		"amount":     amount,
		"currency":   currency,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/internal/payments", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("payment service returned %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var payment CreatePaymentResponse
	if err := json.Unmarshal(respBody, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}
