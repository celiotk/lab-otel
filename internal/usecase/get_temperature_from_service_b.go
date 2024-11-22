package usecase

import (
	"github.com/celiotk/lab-otel/internal/entity"
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
}

func NewTemperatureFromServiceBUsecase(weatherQuery entity.TempFromServiceBInterface) *TemperatureFromServiceBUsecase {
	return &TemperatureFromServiceBUsecase{
		temperatureProvider: weatherQuery,
	}
}

func (u *TemperatureFromServiceBUsecase) Execute(input TemperatureFromServiceBInput) (*TemperatureFromServiceBOutput, error) {
	weather, err := u.temperatureProvider.Get(input.CEP)
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
