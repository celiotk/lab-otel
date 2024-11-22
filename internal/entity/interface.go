package entity

type TemperatureProviderInterface interface {
	Get(city string) (*Weather, error)
}

type LocationProviderInterface interface {
	Get(cep string) (*CEP, error)
}

type TempFromServiceBInterface interface {
	Get(cep string) (*Weather, error)
}
