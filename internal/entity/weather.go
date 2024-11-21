package entity

type Weather struct {
	City           string
	TempCelsius    float64
	TempFahrenheit float64
	TempKelvin     float64
}

func NewWeather(city string, tempCelsius float64) *Weather {
	return &Weather{
		City:        city,
		TempCelsius: tempCelsius,
	}
}

func (w *Weather) ConvertToFahrenheitAndKelvin() {
	w.TempFahrenheit = w.TempCelsius*1.8 + 32
	w.TempKelvin = w.TempCelsius + 273
}
