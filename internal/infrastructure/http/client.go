package http

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type HTTPOrderIntegration struct {
	targetURL  string
	httpClient *http.Client
}

func NewHTTPOrderIntegration(targetURL string) *HTTPOrderIntegration {
	return &HTTPOrderIntegration{
		targetURL: targetURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (h *HTTPOrderIntegration) UpdateOrder(data []byte) error {
	req, err := http.NewRequest(http.MethodPatch, h.targetURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("received non-success status code: %d", resp.StatusCode)
	}

	return nil
}
