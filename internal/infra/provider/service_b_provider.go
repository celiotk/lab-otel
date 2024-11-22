package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/celiotk/lab-otel/internal/entity"
)

type serviceBResult struct {
	City           string  `json:"city"`
	TempCelsius    float64 `json:"temp_C"`
	TempFahrenheit float64 `json:"temp_F"`
	TempKelvin     float64 `json:"temp_K"`
}

type serviceBProvider struct {
	providerAddress string
}

func NewServiceBProvider(addr string) *serviceBProvider {
	return &serviceBProvider{addr}
}

func (c *serviceBProvider) Get(cep string) (*entity.Weather, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/temperature/%s", c.providerAddress, cep), nil)
	if err != nil {
		return nil, fmt.Errorf("serviceBProvider.Get: %w", err)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("serviceBProvider.Get: %w", err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, entity.ErrCepNotFound
	case http.StatusUnprocessableEntity:
		return nil, entity.ErrInvalidCep
	default:
		return nil, fmt.Errorf("serviceBProvider.Get: invalid status code: %d", resp.StatusCode)
	}
	var result serviceBResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("serviceBProvider.Get: %w", err)
	}
	return &entity.Weather{
		City:           result.City,
		TempCelsius:    result.TempCelsius,
		TempFahrenheit: result.TempFahrenheit,
		TempKelvin:     result.TempKelvin,
	}, nil

}
