package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/celiotk/lab-otel/internal/usecase"

	cepValidator "github.com/celiotk/lab-deploy-cloud-run/pkg/cep_validator"
)

type serviceAHandler struct {
	weatherUsecase usecase.TemperatureFromServiceBUsecase
}

func NewServiceAHandler(weatherUsecase usecase.TemperatureFromServiceBUsecase) *serviceAHandler {
	return &serviceAHandler{
		weatherUsecase: weatherUsecase,
	}
}

func (h *serviceAHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	param := struct {
		CEP string `json:"cep"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cep := param.CEP
	if !cepValidator.IsValid(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	input := usecase.TemperatureFromServiceBInput{
		CEP: cep,
	}
	output, err := h.weatherUsecase.Execute(input)
	if err != nil {
		if err == usecase.ErrCepNotFound {
			fmt.Printf("zipcode not found: %s\n", cep)
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		fmt.Printf("error getting temperature: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
