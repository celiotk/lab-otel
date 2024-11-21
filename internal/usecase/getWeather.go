package usecase

import (
	"fmt"

	"github.com/celiotk/lab-otel/internal/entity"
)

var ErrCepNotFound = entity.ErrCepNotFound

type TemperatureFromCepInput struct {
	CEP string
}

type TemperatureFromCepOutput struct {
	Temp_C string `json:"temp_C"`
	Temp_F string `json:"temp_F"`
	Temp_K string `json:"temp_K"`
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
		Temp_C: fmt.Sprintf("%.1f", weather.TempCelsius),
		Temp_F: fmt.Sprintf("%.1f", weather.TempFahrenheit),
		Temp_K: fmt.Sprintf("%.1f", weather.TempKelvin),
	}, nil
}
