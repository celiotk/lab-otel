package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/celiotk/lab-otel/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	cepValidator "github.com/celiotk/lab-deploy-cloud-run/pkg/cep_validator"
)

type temperatureFromServiceBHandler struct {
	weatherUsecase usecase.TemperatureFromServiceBUsecase
	otelTrace      trace.Tracer
}

func NewServiceAHandler(uc usecase.TemperatureFromServiceBUsecase, otelTrace trace.Tracer) *temperatureFromServiceBHandler {
	return &temperatureFromServiceBHandler{
		weatherUsecase: uc,
		otelTrace:      otelTrace,
	}
}

func (h *temperatureFromServiceBHandler) PostTemperature(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)
	ctx, span := h.otelTrace.Start(ctx, "temeratureFromServiceBHandler.PostTemperature")
	defer span.End()

	param := struct {
		CEP string `json:"cep"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.SetStatus(codes.Error, "json decode failed")
		span.RecordError(err)
		return
	}
	cep := param.CEP
	if !cepValidator.IsValid(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		span.SetStatus(codes.Error, "invalid zipcode")
		span.RecordError(errors.New("invalid zipcode"))
		return
	}
	input := usecase.TemperatureFromServiceBInput{
		CEP: cep,
	}
	output, err := h.weatherUsecase.Execute(ctx, input)
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
		span.SetStatus(codes.Error, "json encode failed")
		span.RecordError(err)
		return
	}
}
