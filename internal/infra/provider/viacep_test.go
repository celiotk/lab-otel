package provider

import (
	"testing"
)

func TestGetCep(t *testing.T) {
	cep := NewViaCepProvider()
	result, err := cep.Get("01001-000")
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

	_, err = cep.Get("99999999")
	if err == nil {
		t.Error("Expected error")
	}
}
