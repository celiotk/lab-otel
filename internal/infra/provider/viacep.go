package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/celiotk/lab-otel/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type viaCepResult struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`

	Erro string `json:"erro"`
}

type viaCepProvider struct {
	otelTracer trace.Tracer
}

func NewViaCepProvider(otelTrace trace.Tracer) *viaCepProvider {
	return &viaCepProvider{otelTrace}
}

func (c *viaCepProvider) Get(ctx context.Context, cep string) (*entity.CEP, error) {
	ctx, span := c.otelTracer.Start(ctx, "viaCepProvider.Get")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep), nil)
	if err != nil {
		return nil, fmt.Errorf("viaCepProvider.Get: %w", err)
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, "client.Do failed")
		span.RecordError(err)
		return nil, fmt.Errorf("viaCepProvider.Get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, "client.Do failed")
		span.RecordError(errors.New("invalid status code"))
		return nil, fmt.Errorf("viaCepProvider.Get: invalid status code: %d", resp.StatusCode)
	}
	var cepResult viaCepResult
	if err := json.NewDecoder(resp.Body).Decode(&cepResult); err != nil {
		return nil, fmt.Errorf("viaCepProvider.Get: %w", err)
	}
	if cepResult.Erro != "" {
		span.SetStatus(codes.Error, "client.Do cep not found")
		span.RecordError(fmt.Errorf("cep not found"))
		return nil, entity.ErrCepNotFound
	}
	return &entity.CEP{
		CEP:  cepResult.Cep,
		City: cepResult.Localidade,
	}, nil

}
