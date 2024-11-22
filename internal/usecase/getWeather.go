package usecase

import (
	"github.com/celiotk/lab-otel/internal/entity"
)

var ErrCepNotFound = entity.ErrCepNotFound
var ErrInvalidCep = entity.ErrInvalidCep

type TemperatureFromCepInput struct {
	CEP string
}

type TemperatureFromCepOutput struct {
	City   string  `json:"city"`
	Temp_C float64 `json:"temp_C"`
	Temp_F float64 `json:"temp_F"`
	Temp_K float64 `json:"temp_K"`
}

type TemperatureFromCepUsecase struct {
	temperatureProvider entity.TemperatureProviderInterface
	locationProvider    entity.LocationProviderInterface
}

func NewTemperatureFromCepUsecase(weatherQuery entity.TemperatureProviderInterface, cepQuery entity.LocationProviderInterface) *TemperatureFromCepUsecase {
	return &TemperatureFromCepUsecase{
		temperatureProvider: weatherQuery,
		locationProvider:    cepQuery,
	}
}

func (u *TemperatureFromCepUsecase) Execute(input TemperatureFromCepInput) (*TemperatureFromCepOutput, error) {
	cep, err := u.locationProvider.Get(input.CEP)
	if err != nil {
		if err == entity.ErrCepNotFound {
			return nil, ErrCepNotFound
		}
		return nil, err
	}

	weather, err := u.temperatureProvider.Get(cep.City)
	if err != nil {
		return nil, err
	}

	weather.ConvertToFahrenheitAndKelvin()

	return &TemperatureFromCepOutput{
		City:   cep.City,
		Temp_C: weather.TempCelsius,
		Temp_F: weather.TempFahrenheit,
		Temp_K: weather.TempKelvin,
	}, nil
}
