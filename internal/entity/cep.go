package entity

import "errors"

var ErrCepNotFound = errors.New("CEP not found")
var ErrInvalidCep = errors.New("Invalid CEP")

type CEP struct {
	CEP  string `json:"cep"`
	City string `json:"city"`
}
