package usecase

import (
	"context"

	"github.com/celiotk/lab-otel/internal/entity"
	"go.opentelemetry.io/otel/trace"
)

var ErrCepNotFound = entity.ErrCepNotFound
var ErrInvalidCep = entity.ErrInvalidCep

type TemperatureByCepInput struct {
	CEP string
}

type TemperatureByCepOutput struct {
	City   string  `json:"city"`
	Temp_C float64 `json:"temp_C"`
	Temp_F float64 `json:"temp_F"`
	Temp_K float64 `json:"temp_K"`
}

type TemperatureByCepUsecase struct {
	temperatureProvider entity.TemperatureProviderInterface
	locationProvider    entity.LocationProviderInterface
	otelTracer          trace.Tracer
}

func NewTemperatureByCepUsecase(weatherQuery entity.TemperatureProviderInterface, cepQuery entity.LocationProviderInterface, otelTracer trace.Tracer) *TemperatureByCepUsecase {
	return &TemperatureByCepUsecase{
		temperatureProvider: weatherQuery,
		locationProvider:    cepQuery,
		otelTracer:          otelTracer,
	}
}

func (u *TemperatureByCepUsecase) Execute(ctx context.Context, input TemperatureByCepInput) (*TemperatureByCepOutput, error) {
	ctx, span := u.otelTracer.Start(ctx, "TemperatureByCepUsecase.Execute")
	defer span.End()

	cep, err := u.locationProvider.Get(ctx, input.CEP)
	if err != nil {
		if err == entity.ErrCepNotFound {
			return nil, ErrCepNotFound
		}
		return nil, err
	}

	weather, err := u.temperatureProvider.Get(ctx, cep.City)
	if err != nil {
		return nil, err
	}

	weather.ConvertToFahrenheitAndKelvin()

	return &TemperatureByCepOutput{
		City:   cep.City,
		Temp_C: weather.TempCelsius,
		Temp_F: weather.TempFahrenheit,
		Temp_K: weather.TempKelvin,
	}, nil
}
