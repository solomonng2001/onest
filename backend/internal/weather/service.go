package weather

import (
	"context"
	"sort"
)

type Service struct {
	client *Client
}

func NewService(client *Client) *Service {
	return &Service{
		client: client,
	}
}

// GetForecasts returns the shared valid period and every location's forecast.
func (s *Service) GetForecasts(
	ctx context.Context,
) (*ForecastResponse, error) {
	apiResponse, err := s.client.GetForecast(ctx)
	if err != nil {
		return nil, err
	}

	latest := apiResponse.Data.Items[0]

	forecasts := make(
		[]LocationForecast,
		0,
		len(latest.Forecasts),
	)

	for _, areaForecast := range latest.Forecasts {
		forecasts = append(
			forecasts,
			LocationForecast{
				Location: areaForecast.Area,
				Forecast: areaForecast.Forecast,
			},
		)
	}

	sort.Slice(forecasts, func(i, j int) bool {
		return forecasts[i].Location <
			forecasts[j].Location
	})

	return &ForecastResponse{
		ValidPeriod: latest.ValidPeriod.Text,
		Forecasts:   forecasts,
	}, nil
}
