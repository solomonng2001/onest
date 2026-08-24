package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api-open.data.gov.sg/v2/real-time/api/two-hr-forecast"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetForecast(
	ctx context.Context,
) (*APIResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create weather request: %w",
			err,
		)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"request weather forecast: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"weather API returned status %d",
			resp.StatusCode,
		)
	}

	var result APIResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(
			"decode weather response: %w",
			err,
		)
	}

	if result.Code != 0 {
		if result.ErrorMsg != "" {
			return nil, fmt.Errorf(
				"weather API error: %s",
				result.ErrorMsg,
			)
		}

		return nil, errors.New(
			"weather API returned an unsuccessful response",
		)
	}

	if len(result.Data.Items) == 0 {
		return nil, errors.New(
			"weather API returned no forecast items",
		)
	}

	if len(result.Data.Items[0].Forecasts) == 0 {
		return nil, errors.New(
			"weather API returned no location forecasts",
		)
	}

	return &result, nil
}
