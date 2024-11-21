package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/celiotk/lab-otel/internal/entity"
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
}

func NewViaCepProvider() *viaCepProvider {
	return &viaCepProvider{}
}

func (c *viaCepProvider) Get(cep string) (*entity.CEP, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep), nil)
	if err != nil {
		return nil, fmt.Errorf("viaCepProvider.Get: %w", err)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("viaCepProvider.Get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("viaCepProvider.Get: invalid status code: %d", resp.StatusCode)
	}
	var cepResult viaCepResult
	if err := json.NewDecoder(resp.Body).Decode(&cepResult); err != nil {
		return nil, fmt.Errorf("viaCepProvider.Get: %w", err)
	}
	if cepResult.Erro != "" {
		return nil, entity.ErrCepNotFound
	}
	return &entity.CEP{
		CEP:  cepResult.Cep,
		City: cepResult.Localidade,
	}, nil

}
