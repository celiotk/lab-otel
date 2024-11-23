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
	"github.com/go-chi/chi/v5"
)

type temperatureByCepHandler struct {
	weatherUsecase usecase.TemperatureByCepUsecase
	otelTracer     trace.Tracer
}

func NewTemperatureByCepHandler(uc usecase.TemperatureByCepUsecase, otelTracer trace.Tracer) *temperatureByCepHandler {
	return &temperatureByCepHandler{
		weatherUsecase: uc,
		otelTracer:     otelTracer,
	}
}

func (h *temperatureByCepHandler) GetTemperature(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)
	ctx, span := h.otelTracer.Start(ctx, "temperatureByCepHandler.GetTemperature")
	defer span.End()

	cep := chi.URLParam(r, "cep")
	if !cepValidator.IsValid(cep) {
		fmt.Printf("invalid zipcode: %s\n", cep)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		span.SetStatus(codes.Error, "invalid zipcode")
		span.RecordError(errors.New("invalid zipcode"))
		return
	}
	input := usecase.TemperatureByCepInput{
		CEP: cep,
	}
	output, err := h.weatherUsecase.Execute(ctx, input)
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
		span.SetStatus(codes.Error, "json.NewEncoder.Encode failed")
		span.RecordError(err)
		return
	}
}
