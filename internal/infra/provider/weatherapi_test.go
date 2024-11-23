package provider

import (
	"context"
	"testing"

	"github.com/celiotk/lab-otel/configs"
	"go.opentelemetry.io/otel"
)

func TestWeatherApiProvider(t *testing.T) {
	cfg, err := configs.LoadConfig("../../../")
	if err != nil {
		t.Error(err)
		return
	}
	weather := NewWeatherApiProvider(cfg.WeatherAPIKey, otel.Tracer("test"))
	result, err := weather.Get(context.Background(), "São Paulo")
	if err != nil {
		t.Error(err)
		return
	}
	if result.City != "São Paulo" {
		t.Errorf("City not found: expected %s, got %s", "São Paulo", result.City)
	}
}
