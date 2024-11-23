package provider

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestGetCep(t *testing.T) {
	cep := NewViaCepProvider(otel.Tracer("test"))
	result, err := cep.Get(context.Background(), "01001-000")
	if err != nil {
		t.Error(err)
		return
	}
	if result.CEP != "01001-000" {
		t.Error("CEP not found")
	}
	if result.City != "São Paulo" {
		t.Error("City not found")
	}

	_, err = cep.Get(context.Background(), "99999999")
	if err == nil {
		t.Error("Expected error")
	}
}
