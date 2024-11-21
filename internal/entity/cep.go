package entity

import "errors"

var ErrCepNotFound = errors.New("CEP not found")

type CEP struct {
	CEP  string `json:"cep"`
	City string `json:"city"`
}
