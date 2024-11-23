package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/celiotk/lab-otel/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type serviceBResult struct {
	City           string  `json:"city"`
	TempCelsius    float64 `json:"temp_C"`
	TempFahrenheit float64 `json:"temp_F"`
	TempKelvin     float64 `json:"temp_K"`
}

type serviceBProvider struct {
	providerAddress string
	otelTracer      trace.Tracer
}

func NewServiceBProvider(addr string, otelTracer trace.Tracer) *serviceBProvider {
	return &serviceBProvider{addr, otelTracer}
}

func (c *serviceBProvider) Get(ctx context.Context, cep string) (*entity.Weather, error) {
	ctx, span := c.otelTracer.Start(ctx, "serviceBProvider.Get")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/temperature/%s", c.providerAddress, cep), nil)
	if err != nil {
		return nil, fmt.Errorf("serviceBProvider.Get: %w", err)
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, "client.Do failed")
		span.RecordError(err)
		return nil, fmt.Errorf("serviceBProvider.Get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, "client.Do failed")
		span.RecordError(errors.New("invalid status code"))
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, entity.ErrCepNotFound
		case http.StatusUnprocessableEntity:
			return nil, entity.ErrInvalidCep
		default:
			return nil, fmt.Errorf("serviceBProvider.Get: invalid status code: %d", resp.StatusCode)
		}
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
