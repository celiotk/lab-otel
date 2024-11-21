package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/celiotk/lab-otel/internal/usecase"

	cepValidator "github.com/celiotk/lab-otel/pkg/cep_validator"
	"github.com/go-chi/chi/v5"
)

type WeatherHandler struct {
	weatherUsecase usecase.TemperatureFromCepUsecase
}

func NewWeatherHandler(weatherUsecase usecase.TemperatureFromCepUsecase) *WeatherHandler {
	return &WeatherHandler{
		weatherUsecase: weatherUsecase,
	}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if !cepValidator.IsValid(cep) {
		fmt.Printf("invalid zipcode: %s\n", cep)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	input := usecase.TemperatureFromCepInput{
		CEP: cep,
	}
	output, err := h.weatherUsecase.Execute(input)
	if err != nil {
		if err == usecase.ErrCepNotFound {
			fmt.Printf("zipcode not found: %s\n", cep)
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		fmt.Printf("error getting weather: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
