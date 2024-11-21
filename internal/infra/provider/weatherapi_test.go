package provider

import (
	"context"
	"testing"

	"github.com/celiotk/lab-otel/configs"
)

func TestWeatherApiProvider(t *testing.T) {
	cfg, err := configs.LoadConfig("../../../")
	if err != nil {
		t.Error(err)
		return
	}
	weather := NewWeatherApiProvider(context.Background(), cfg.WeatherAPIKey)
	result, err := weather.Get("São Paulo")
	if err != nil {
		t.Error(err)
		return
	}
	if result.City != "São Paulo" {
		t.Errorf("City not found: expected %s, got %s", "São Paulo", result.City)
	}
}
