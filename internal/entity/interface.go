package entity

import "context"

type TemperatureProviderInterface interface {
	Get(ctx context.Context, city string) (*Weather, error)
}

type LocationProviderInterface interface {
	Get(ctx context.Context, cep string) (*CEP, error)
}

type TempFromServiceBInterface interface {
	Get(ctx context.Context, cep string) (*Weather, error)
}
