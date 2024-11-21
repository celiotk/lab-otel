package cepValidator

import "regexp"

func IsValid(cep string) bool {
	regexp := regexp.MustCompile(`^\d{5}-?\d{3}$`)
	return regexp.MatchString(cep)
}
