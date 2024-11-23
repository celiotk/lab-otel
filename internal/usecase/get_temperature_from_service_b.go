package usecase

import (
	"context"

	"github.com/celiotk/lab-otel/internal/entity"
	"go.opentelemetry.io/otel/trace"
)

type TemperatureFromServiceBInput struct {
	CEP string
}

type TemperatureFromServiceBOutput struct {
	City   string  `json:"city"`
	Temp_C float64 `json:"temp_C"`
	Temp_F float64 `json:"temp_F"`
	Temp_K float64 `json:"temp_K"`
}

type TemperatureFromServiceBUsecase struct {
	temperatureProvider entity.TempFromServiceBInterface
	otelTracer          trace.Tracer
}

func NewTemperatureFromServiceBUsecase(weatherQuery entity.TempFromServiceBInterface, otelTracer trace.Tracer) *TemperatureFromServiceBUsecase {
	return &TemperatureFromServiceBUsecase{
		temperatureProvider: weatherQuery,
		otelTracer:          otelTracer,
	}
}

func (u *TemperatureFromServiceBUsecase) Execute(ctx context.Context, input TemperatureFromServiceBInput) (*TemperatureFromServiceBOutput, error) {
	ctx, span := u.otelTracer.Start(ctx, "TemperatureFromServiceBUsecase.Execute")
	defer span.End()

	weather, err := u.temperatureProvider.Get(ctx, input.CEP)
	switch err {
	case entity.ErrCepNotFound:
		return nil, ErrCepNotFound
	case entity.ErrInvalidCep:
		return nil, ErrInvalidCep
	case nil:
	default:
		return nil, err
	}

	return &TemperatureFromServiceBOutput{
		City:   weather.City,
		Temp_C: weather.TempCelsius,
		Temp_F: weather.TempFahrenheit,
		Temp_K: weather.TempKelvin,
	}, nil
}
