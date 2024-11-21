package cepValidator

import (
	"testing"
)

func TestValidAndInvalidCep(t *testing.T) {
	zipCodeValidationMap := map[string]bool{
		"01001-000": true,
		"01001000":  true,
		"01001-00":  false,
		"01001-ABC": false,
		"01001ABC":  false,
	}
	for cep, expected := range zipCodeValidationMap {
		if IsValid(cep) != expected {
			t.Errorf("CEP %s is invalid", cep)
		}
	}
}
